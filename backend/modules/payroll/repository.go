package payroll

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

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

// ErrNotEligible is returned when salary slips are disabled for the tenant.
var ErrNotEligible = errors.New("payroll: salary slips are disabled for this tenant")

// ErrNoStructure is returned when the employee has no salary structure.
var ErrNoStructure = errors.New("payroll: no salary structure for this employee")
