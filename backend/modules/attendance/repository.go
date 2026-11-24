package attendance

import (
	"context"
	"time"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/db"
)

// Repository is the tenant-scoped pgx connection to the attendance schema.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository returns the attendance repository backed by pool.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// tenantTx acquires a tenant-scoped connection for tenantID.
func (r *Repository) tenantTx(ctx context.Context, tenantID string) (*pgxpool.Conn, func(), error) {
	return db.NewTenantTx(ctx, r.pool, tenantID)
}

// scanPolicy reads one attendance.tenant_policies row, returning
// ErrPolicyNotFound when unset.
func (r *Repository) scanPolicy(ctx context.Context, tenantID string, p *TenantPolicy) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	rows, err := conn.Query(ctx,
		`SELECT methods, shift_start_time, shift_end_time
		 FROM attendance.tenant_policies WHERE tenant_id = $1`, tenantID)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return ErrPolicyNotFound{tenantID: tenantID}
	}
	if err := rows.Scan(&p.Methods, &p.ShiftStart, &p.ShiftEnd); err != nil {
		return err
	}
	return nil
}

// GetLogs returns the tenant's raw punch events.
func (r *Repository) GetLogs(ctx context.Context, tenantID string) ([]PunchLog, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	rows, err := conn.Query(ctx,
		`SELECT id, tenant_id, employee_id, method, start_time, end_time
		 FROM attendance.logs WHERE tenant_id = $1
		 ORDER BY start_time DESC LIMIT 1000`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PunchLog
	for rows.Next() {
		var l PunchLog
		if err := rows.Scan(&l.ID, &l.TenantID, &l.EmployeeID, &l.Method, &l.Start, &l.End); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

// GetTenantPolicy returns the tenant's locally-cached policy.
func (r *Repository) GetTenantPolicy(ctx context.Context, tenantID string) (TenantPolicy, error) {
	var p TenantPolicy
	if err := r.scanPolicy(ctx, tenantID, &p); err != nil {
		return p, err
	}
	return p, nil
}

// SavePolicy writes the tenant's locally-cached policy.
func (r *Repository) SavePolicy(ctx context.Context, tenantID string, p TenantPolicy) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`INSERT INTO attendance.tenant_policies (tenant_id, methods, shift_start_time, shift_end_time)
		 VALUES ($1, $2::punch_method[], $3, $4)
		 ON CONFLICT (tenant_id) DO UPDATE SET
		 methods = EXCLUDED.methods,
		 shift_start_time = EXCLUDED.shift_start_time,
		 shift_end_time = EXCLUDED.shift_end_time`,
		tenantID, p.Methods, p.ShiftStart, p.ShiftEnd)
	return err
}

// InsertLog records a raw punch event.
func (r *Repository) InsertLog(ctx context.Context, tenantID, employeeID string, method PunchMethod, ts time.Time) (int64, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return 0, err
	}
	defer release()

	var id int64
	err = conn.QueryRow(ctx,
		`INSERT INTO attendance.logs (tenant_id, employee_id, method, start_time, end_time)
		 VALUES ($1, $2, $3, $4::time, $5::time) RETURNING id`,
		tenantID, employeeID, method, ts, ts).Scan(&id)
	return id, err
}

// CreateShift creates a new shift for a tenant
func (r *Repository) CreateShift(ctx context.Context, tenantID string, shift Shift) (*Shift, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var createdShift Shift
	err = conn.QueryRow(ctx,
		`INSERT INTO attendance.shifts (id, tenant_id, name, description, start_time, end_time, is_active)
		 VALUES ($1, $2, $3, $4, $5::time, $6::time, $7) RETURNING id, tenant_id, name, description, start_time, end_time, is_active, created_at, updated_at`,
		uuid.New().String(), tenantID, shift.Name, shift.Description, shift.StartTime, shift.EndTime, shift.IsActive).Scan(
		&createdShift.ID, &createdShift.TenantID, &createdShift.Name, &createdShift.Description,
		&createdShift.StartTime, &createdShift.EndTime, &createdShift.IsActive,
		&createdShift.CreatedAt, &createdShift.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &createdShift, nil
}

// GetShift retrieves a shift by ID
func (r *Repository) GetShift(ctx context.Context, tenantID string, shiftID string) (*Shift, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var shift Shift
	err = conn.QueryRow(ctx,
		`SELECT id, tenant_id, name, description, start_time, end_time, is_active, created_at, updated_at
		 FROM attendance.shifts WHERE id = $1 AND tenant_id = $2`,
		shiftID, tenantID).Scan(
		&shift.ID, &shift.TenantID, &shift.Name, &shift.Description,
		&shift.StartTime, &shift.EndTime, &shift.IsActive,
		&shift.CreatedAt, &shift.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &shift, nil
}

// ListShifts lists shifts for a tenant
func (r *Repository) ListShifts(ctx context.Context, tenantID string, page, limit int) ([]*Shift, int, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	offset := (page - 1) * limit

	rows, err := conn.Query(ctx,
		`SELECT id, tenant_id, name, description, start_time, end_time, is_active, created_at, updated_at
		 FROM attendance.shifts WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shifts []*Shift
	for rows.Next() {
		var shift Shift
		err = rows.Scan(&shift.ID, &shift.TenantID, &shift.Name, &shift.Description,
			&shift.StartTime, &shift.EndTime, &shift.IsActive,
			&shift.CreatedAt, &shift.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		shifts = append(shifts, &shift)
	}

	// Get total count
	var total int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM attendance.shifts WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

// UpdateShift updates an existing shift
func (r *Repository) UpdateShift(ctx context.Context, tenantID string, shiftID string, shift Shift) (*Shift, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var updatedShift Shift
	err = conn.QueryRow(ctx,
		`UPDATE attendance.shifts SET name = COALESCE($1, name),
		 description = COALESCE($2, description),
		 start_time = COALESCE($3, start_time),
		 end_time = COALESCE($4, end_time),
		 is_active = COALESCE($5, is_active),
		 updated_at = NOW()
		 WHERE id = $6 AND tenant_id = $7 RETURNING id, tenant_id, name, description, start_time, end_time, is_active, created_at, updated_at`,
		shift.Name, shift.Description, shift.StartTime, shift.EndTime, shift.IsActive,
		shiftID, tenantID).Scan(
		&updatedShift.ID, &updatedShift.TenantID, &updatedShift.Name, &updatedShift.Description,
		&updatedShift.StartTime, &updatedShift.EndTime, &updatedShift.IsActive,
		&updatedShift.CreatedAt, &updatedShift.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedShift, nil
}

// DeleteShift deletes a shift
func (r *Repository) DeleteShift(ctx context.Context, tenantID string, shiftID string) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM attendance.shifts WHERE id = $1 AND tenant_id = $2`, shiftID, tenantID)
	return err
}

// CreateAttendanceRule creates a new attendance rule for a tenant
func (r *Repository) CreateAttendanceRule(ctx context.Context, tenantID string, rule AttendanceRule) (*AttendanceRule, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var createdRule AttendanceRule
	err = conn.QueryRow(ctx,
		`INSERT INTO attendance.attendance_rules (id, tenant_id, rule_name, description, rule_type, condition, action, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, tenant_id, rule_name, description, rule_type, condition, action, is_active, created_at, updated_at`,
		uuid.New().String(), tenantID, rule.Name, rule.Description, rule.RuleType,
		rule.Condition, rule.Action, rule.IsActive).Scan(
		&createdRule.ID, &createdRule.TenantID, &createdRule.Name, &createdRule.Description,
		&createdRule.RuleType, &createdRule.Condition, &createdRule.Action, &createdRule.IsActive,
		&createdRule.CreatedAt, &createdRule.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &createdRule, nil
}

// GetAttendanceRule retrieves an attendance rule by ID
func (r *Repository) GetAttendanceRule(ctx context.Context, tenantID string, ruleID string) (*AttendanceRule, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var rule AttendanceRule
	err = conn.QueryRow(ctx,
		`SELECT id, tenant_id, rule_name, description, rule_type, condition, action, is_active, created_at, updated_at
		 FROM attendance.attendance_rules WHERE id = $1 AND tenant_id = $2`,
		ruleID, tenantID).Scan(
		&rule.ID, &rule.TenantID, &rule.Name, &rule.Description,
		&rule.RuleType, &rule.Condition, &rule.Action, &rule.IsActive,
		&rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

// ListAttendanceRules lists attendance rules for a tenant
func (r *Repository) ListAttendanceRules(ctx context.Context, tenantID string, page, limit int) ([]*AttendanceRule, int, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	offset := (page - 1) * limit

	rows, err := conn.Query(ctx,
		`SELECT id, tenant_id, rule_name, description, rule_type, condition, action, is_active, created_at, updated_at
		 FROM attendance.attendance_rules WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rules []*AttendanceRule
	for rows.Next() {
		var rule AttendanceRule
		err = rows.Scan(&rule.ID, &rule.TenantID, &rule.Name, &rule.Description,
			&rule.RuleType, &rule.Condition, &rule.Action, &rule.IsActive,
			&rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, &rule)
	}

	// Get total count
	var total int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM attendance.attendance_rules WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

// UpdateAttendanceRule updates an existing attendance rule
func (r *Repository) UpdateAttendanceRule(ctx context.Context, tenantID string, ruleID string, rule AttendanceRule) (*AttendanceRule, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var updatedRule AttendanceRule
	err = conn.QueryRow(ctx,
		`UPDATE attendance.attendance_rules SET rule_name = COALESCE($1, rule_name),
		 description = COALESCE($2, description),
		 rule_type = COALESCE($3, rule_type),
		 condition = COALESCE($4, condition),
		 action = COALESCE($5, action),
		 is_active = COALESCE($6, is_active),
		 updated_at = NOW()
		 WHERE id = $7 AND tenant_id = $8 RETURNING id, tenant_id, rule_name, description, rule_type, condition, action, is_active, created_at, updated_at`,
		rule.Name, rule.Description, rule.RuleType, rule.Condition, rule.Action, rule.IsActive,
		ruleID, tenantID).Scan(
		&updatedRule.ID, &updatedRule.TenantID, &updatedRule.Name, &updatedRule.Description,
		&updatedRule.RuleType, &updatedRule.Condition, &updatedRule.Action, &updatedRule.IsActive,
		&updatedRule.CreatedAt, &updatedRule.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedRule, nil
}

// DeleteAttendanceRule deletes an attendance rule
func (r *Repository) DeleteAttendanceRule(ctx context.Context, tenantID string, ruleID string) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM attendance.attendance_rules WHERE id = $1 AND tenant_id = $2`, ruleID, tenantID)
	return err
}
