package organization

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/db"
)

// ErrPolicyNotFound is returned when a tenant has no policy yet.
type ErrPolicyNotFound struct{ tenantID string }

func (e ErrPolicyNotFound) Error() string {
	return "org: policy not found for tenant " + e.tenantID
}

// ErrNotFound is the zero-value ErrPolicyNotFound.
var ErrNotFound = ErrPolicyNotFound{}

// Repository is the tenant-scoped pgx connection to the org schema.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository returns the org repository backed by pool.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// tenantTx acquires a tenant-scoped connection for tenantID.
func (r *Repository) tenantTx(ctx context.Context, tenantID string) (*pgxpool.Conn, func(), error) {
	return db.NewTenantTx(ctx, r.pool, tenantID)
}

// GetPolicy returns the tenant's attendance policy.
func (r *Repository) GetPolicy(ctx context.Context, tenantID string) (OrgPolicy, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return OrgPolicy{}, err
	}
	defer release()

	var p OrgPolicy
	rows, err := conn.Query(ctx,
		`SELECT methods, salary_slip_enabled, shift_start_time, shift_end_time
		 FROM org.policies WHERE tenant_id = $1`, tenantID)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	if !rows.Next() {
		return OrgPolicy{}, ErrPolicyNotFound{tenantID: tenantID}
	}
	if err := rows.Scan(
		&p.Methods, &p.SalarySlipEnabled, &p.ShiftStart, &p.ShiftEnd); err != nil {
		return p, err
	}
	return p, nil
}

// SavePolicy writes the tenant's attendance policy.
func (r *Repository) SavePolicy(ctx context.Context, tenantID string, p OrgPolicy) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`INSERT INTO org.policies (tenant_id, methods, salary_slip_enabled, shift_start_time, shift_end_time)
		 VALUES ($1, $2::text[], $3, $4, $5)
		 ON CONFLICT (tenant_id) DO UPDATE SET
		 methods = EXCLUDED.methods,
		 salary_slip_enabled = EXCLUDED.salary_slip_enabled,
		 shift_start_time = EXCLUDED.shift_start_time,
		 shift_end_time = EXCLUDED.shift_end_time`,
		tenantID, p.Methods, p.SalarySlipEnabled, p.ShiftStart, p.ShiftEnd)
	return err
}

// GetEmployees lists the tenant's Keycloak-backed users.
func (r *Repository) GetEmployees(ctx context.Context, tenantID string) ([]Employee, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	rows, err := conn.Query(ctx,
		`SELECT id, keycloak_sub, role FROM org.employees WHERE tenant_id = $1
		 ORDER BY id`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.KeycloakSub, &e.Role); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}

// CreateEmployee saves a Keycloak user and returns its keycloak_sub.
func (r *Repository) CreateEmployee(ctx context.Context, tenantID, keycloakSub, role string) (string, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return "", err
	}
	defer release()

	var sub string
	err = conn.QueryRow(ctx,
		`INSERT INTO org.employees (tenant_id, keycloak_sub, role) VALUES ($1, $2, $3) RETURNING keycloak_sub`,
		tenantID, keycloakSub, role).Scan(&sub)
	return sub, err
}
