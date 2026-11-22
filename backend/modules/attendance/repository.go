package attendance

import (
	"context"
	"time"

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
