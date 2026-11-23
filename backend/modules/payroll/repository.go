package payroll

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"

	"hunrx.com/attendance_management/pkg/db"
)

// ErrSettingsNotFound is returned when a tenant has no settings yet.
type ErrSettingsNotFound struct{ tenantID string }

func (e ErrSettingsNotFound) Error() string {
	return "payroll: settings not found for tenant " + e.tenantID
}

// Repository is the tenant-scoped pgx connection to the payroll schema.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository returns the payroll repository backed by pool.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// tenantTx acquires a tenant-scoped connection for tenantID.
func (r *Repository) tenantTx(ctx context.Context, tenantID string) (*pgxpool.Conn, func(), error) {
	return db.NewTenantTx(ctx, r.pool, tenantID)
}

// GetTenantSettings returns the tenant's salary slip eligibility.
func (r *Repository) GetTenantSettings(ctx context.Context, tenantID string) (TenantSettings, error) {
	var t TenantSettings
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return t, err
	}
	defer release()

	rows, err := conn.Query(ctx,
		`SELECT salary_slip_enabled FROM payroll.tenant_settings WHERE tenant_id = $1`, tenantID)
	if err != nil {
		return t, err
	}
	defer rows.Close()
	if !rows.Next() {
		return t, ErrSettingsNotFound{tenantID: tenantID}
	}
	if err := rows.Scan(&t.SalarySlipEnabled); err != nil {
		return t, err
	}
	return t, nil
}

// SaveSettings writes the tenant's salary slip eligibility.
func (r *Repository) SaveSettings(ctx context.Context, tenantID string, t TenantSettings) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`INSERT INTO payroll.tenant_settings (tenant_id, salary_slip_enabled)
		 VALUES ($1, $2)
		 ON CONFLICT (tenant_id) DO UPDATE SET
		 salary_slip_enabled = EXCLUDED.salary_slip_enabled`,
		tenantID, t.SalarySlipEnabled)
	return err
}

// GetStructures returns the tenant's salary structures.
func (r *Repository) GetStructures(ctx context.Context, tenantID string) ([]SalaryStructure, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	rows, err := conn.Query(ctx,
		`SELECT employee_id, base_salary, days_per_month
		 FROM payroll.salary_structures WHERE tenant_id = $1
		 ORDER BY employee_id`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SalaryStructure
	for rows.Next() {
		var s SalaryStructure
		if err := rows.Scan(&s.EmployeeID, &s.BaseSalary, &s.DaysPerMonth); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// SaveStructure writes a salary structure.
func (r *Repository) SaveStructure(ctx context.Context, tenantID, employeeID string, s SalaryStructure) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`INSERT INTO payroll.salary_structures (tenant_id, employee_id, base_salary, days_per_month)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (tenant_id, employee_id) DO UPDATE SET
		 base_salary = EXCLUDED.base_salary,
		 days_per_month = EXCLUDED.days_per_month`,
		tenantID, employeeID, s.BaseSalary, s.DaysPerMonth)
	return err
}

// GetSlips returns the tenant's generated salary slips.
func (r *Repository) GetSlips(ctx context.Context, tenantID string) ([]SalarySlip, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	rows, err := conn.Query(ctx,
		`SELECT id, tenant_id, employee_id, period, base_salary, days
		 FROM payroll.salary_slips WHERE tenant_id = $1
		 ORDER BY id DESC LIMIT 100`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SalarySlip
	for rows.Next() {
		var slip SalarySlip
		if err := rows.Scan(&slip.ID, &slip.TenantID, &slip.EmployeeID, &slip.Period, &slip.BaseSalary, &slip.Days); err != nil {
			return nil, err
		}
		out = append(out, slip)
	}
	return out, nil
}

// InsertSlip writes a generated salary slip.
func (r *Repository) InsertSlip(ctx context.Context, tenantID, employeeID string, period string, slip SalarySlip) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`INSERT INTO payroll.salary_slips (tenant_id, employee_id, period, base_salary, days)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (tenant_id, employee_id, period) DO UPDATE SET
		 base_salary = EXCLUDED.base_salary,
		 days = EXCLUDED.days`,
		tenantID, employeeID, period, slip.BaseSalary, slip.Days)
	return err
}

// CreateStatutoryCompliance creates a new statutory compliance rule
func (r *Repository) CreateStatutoryCompliance(ctx context.Context, tenantID string, compliance StatutoryCompliance) (*StatutoryCompliance, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var created StatutoryCompliance
	err = conn.QueryRow(ctx,
		`INSERT INTO payroll.statutory_compliance (id, tenant_id, compliance_type, applicable_from, applicable_to, rate, max_limit, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, tenant_id, compliance_type, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at`,
		uuid.New().String(), tenantID, compliance.ComplianceType, compliance.ApplicableFrom, compliance.ApplicableTo,
		compliance.Rate, compliance.MaxLimit, compliance.IsActive).Scan(
		&created.ID, &created.TenantID, &created.ComplianceType, &created.ApplicableFrom, &created.ApplicableTo,
		&created.Rate, &created.MaxLimit, &created.IsActive, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &created, nil
}

// GetStatutoryCompliance retrieves a statutory compliance rule by ID
func (r *Repository) GetStatutoryCompliance(ctx context.Context, tenantID string, id string) (*StatutoryCompliance, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var compliance StatutoryCompliance
	err = conn.QueryRow(ctx,
		`SELECT id, tenant_id, compliance_type, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at
		 FROM payroll.statutory_compliance WHERE id = $1 AND tenant_id = $2`,
		id, tenantID).Scan(
		&compliance.ID, &compliance.TenantID, &compliance.ComplianceType, &compliance.ApplicableFrom, &compliance.ApplicableTo,
		&compliance.Rate, &compliance.MaxLimit, &compliance.IsActive, &compliance.CreatedAt, &compliance.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &compliance, nil
}

// ListStatutoryCompliance lists statutory compliance rules for a tenant
func (r *Repository) ListStatutoryCompliance(ctx context.Context, tenantID string, page, limit int) ([]*StatutoryCompliance, int, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	offset := (page - 1) * limit
	
	rows, err := conn.Query(ctx,
		`SELECT id, tenant_id, compliance_type, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at
		 FROM payroll.statutory_compliance WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var compliances []*StatutoryCompliance
	for rows.Next() {
		var compliance StatutoryCompliance
		err = rows.Scan(&compliance.ID, &compliance.TenantID, &compliance.ComplianceType, &compliance.ApplicableFrom, &compliance.ApplicableTo,
			&compliance.Rate, &compliance.MaxLimit, &compliance.IsActive, &compliance.CreatedAt, &compliance.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		compliances = append(compliances, &compliance)
	}
	
	// Get total count
	var total int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM payroll.statutory_compliance WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	return compliances, total, nil
}

// UpdateStatutoryCompliance updates an existing statutory compliance rule
func (r *Repository) UpdateStatutoryCompliance(ctx context.Context, tenantID string, id string, compliance StatutoryCompliance) (*StatutoryCompliance, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var updated StatutoryCompliance
	err = conn.QueryRow(ctx,
		`UPDATE payroll.statutory_compliance SET compliance_type = COALESCE($1, compliance_type),
		 applicable_from = COALESCE($2, applicable_from),
		 applicable_to = COALESCE($3, applicable_to),
		 rate = COALESCE($4, rate),
		 max_limit = COALESCE($5, max_limit),
		 is_active = COALESCE($6, is_active),
		 updated_at = NOW()
		 WHERE id = $7 AND tenant_id = $8 RETURNING id, tenant_id, compliance_type, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at`,
		compliance.ComplianceType, compliance.ApplicableFrom, compliance.ApplicableTo,
		compliance.Rate, compliance.MaxLimit, compliance.IsActive,
		id, tenantID).Scan(
		&updated.ID, &updated.TenantID, &updated.ComplianceType, &updated.ApplicableFrom, &updated.ApplicableTo,
		&updated.Rate, &updated.MaxLimit, &updated.IsActive, &updated.CreatedAt, &updated.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &updated, nil
}

// DeleteStatutoryCompliance deletes a statutory compliance rule
func (r *Repository) DeleteStatutoryCompliance(ctx context.Context, tenantID string, id string) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM payroll.statutory_compliance WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// CreateEmployeeStatutoryDeduction creates employee-specific statutory deductions
func (r *Repository) CreateEmployeeStatutoryDeduction(ctx context.Context, employeeID string, deduction EmployeeStatutoryDeduction) (*EmployeeStatutoryDeduction, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var created EmployeeStatutoryDeduction
	err = conn.QueryRow(ctx,
		`INSERT INTO payroll.employee_statutory_deductions (id, employee_id, statutory_id, applicable_from, applicable_to, rate, max_limit, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, employee_id, statutory_id, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at`,
		uuid.New().String(), employeeID, deduction.StatutoryID, deduction.ApplicableFrom, deduction.ApplicableTo,
		deduction.Rate, deduction.MaxLimit, deduction.IsActive).Scan(
		&created.ID, &created.EmployeeID, &created.StatutoryID, &created.ApplicableFrom, &created.ApplicableTo,
		&created.Rate, &created.MaxLimit, &created.IsActive, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &created, nil
}

// GetEmployeeStatutoryDeduction retrieves employee-specific statutory deduction by ID
func (r *Repository) GetEmployeeStatutoryDeduction(ctx context.Context, employeeID string, id string) (*EmployeeStatutoryDeduction, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var deduction EmployeeStatutoryDeduction
	err = conn.QueryRow(ctx,
		`SELECT id, employee_id, statutory_id, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at
		 FROM payroll.employee_statutory_deductions WHERE id = $1 AND employee_id = $2`,
		id, employeeID).Scan(
		&deduction.ID, &deduction.EmployeeID, &deduction.StatutoryID, &deduction.ApplicableFrom, &deduction.ApplicableTo,
		&deduction.Rate, &deduction.MaxLimit, &deduction.IsActive, &deduction.CreatedAt, &deduction.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &deduction, nil
}

// ListEmployeeStatutoryDeductions lists employee-specific statutory deductions
func (r *Repository) ListEmployeeStatutoryDeductions(ctx context.Context, employeeID string, page, limit int) ([]*EmployeeStatutoryDeduction, int, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, 0, err
	}
	defer release()

	offset := (page - 1) * limit
	
	rows, err := conn.Query(ctx,
		`SELECT id, employee_id, statutory_id, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at
		 FROM payroll.employee_statutory_deductions WHERE employee_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		employeeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var deductions []*EmployeeStatutoryDeduction
	for rows.Next() {
		var deduction EmployeeStatutoryDeduction
		err = rows.Scan(&deduction.ID, &deduction.EmployeeID, &deduction.StatutoryID, &deduction.ApplicableFrom, &deduction.ApplicableTo,
			&deduction.Rate, &deduction.MaxLimit, &deduction.IsActive, &deduction.CreatedAt, &deduction.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		deductions = append(deductions, &deduction)
	}
	
	// Get total count
	var total int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM payroll.employee_statutory_deductions WHERE employee_id = $1`, employeeID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	return deductions, total, nil
}

// UpdateEmployeeStatutoryDeduction updates an employee-specific statutory deduction
func (r *Repository) UpdateEmployeeStatutoryDeduction(ctx context.Context, employeeID string, id string, deduction EmployeeStatutoryDeduction) (*EmployeeStatutoryDeduction, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var updated EmployeeStatutoryDeduction
	err = conn.QueryRow(ctx,
		`UPDATE payroll.employee_statutory_deductions SET statutory_id = COALESCE($1, statutory_id),
		 applicable_from = COALESCE($2, applicable_from),
		 applicable_to = COALESCE($3, applicable_to),
		 rate = COALESCE($4, rate),
		 max_limit = COALESCE($5, max_limit),
		 is_active = COALESCE($6, is_active),
		 updated_at = NOW()
		 WHERE id = $7 AND employee_id = $8 RETURNING id, employee_id, statutory_id, applicable_from, applicable_to, rate, max_limit, is_active, created_at, updated_at`,
		deduction.StatutoryID, deduction.ApplicableFrom, deduction.ApplicableTo,
		deduction.Rate, deduction.MaxLimit, deduction.IsActive,
		id, employeeID).Scan(
		&updated.ID, &updated.EmployeeID, &updated.StatutoryID, &updated.ApplicableFrom, &updated.ApplicableTo,
		&updated.Rate, &updated.MaxLimit, &updated.IsActive, &updated.CreatedAt, &updated.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &updated, nil
}

// DeleteEmployeeStatutoryDeduction deletes an employee-specific statutory deduction
func (r *Repository) DeleteEmployeeStatutoryDeduction(ctx context.Context, employeeID string, id string) error {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM payroll.employee_statutory_deductions WHERE id = $1 AND employee_id = $2`, id, employeeID)
	return err
}

// ErrNotEligible is returned when salary slips are disabled for the tenant.
var ErrNotEligible = errors.New("payroll: salary slips are disabled for this tenant")

// ErrNoStructure is returned when the employee has no salary structure.
var ErrNoStructure = errors.New("payroll: no salary structure for this employee")
