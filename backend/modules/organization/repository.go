package organization

import (
	"context"
	"time"

	"github.com/google/uuid"
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

// CreateOrganization creates a new organization.
func (r *Repository) CreateOrganization(ctx context.Context, id uuid.UUID, req CreateOrganizationRequest) (*Organization, error) {
	conn, release, err := r.tenantTx(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var org Organization
	err = conn.QueryRow(ctx,
		`INSERT INTO org.organizations (id, name, description, industry, website, logo_url, timezone, currency, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, name, description, industry, website, logo_url, timezone, currency, created_at`,
		id, req.Name, req.Description, req.Industry, req.Website, req.LogoURL, req.Timezone, req.Currency, req.CreatedAt).Scan(
		&org.ID, &org.Name, &org.Description, &org.Industry, &org.Website, &org.LogoURL, &org.Timezone, &org.Currency, &org.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &org, nil
}

// GetOrganization retrieves an organization by ID.
func (r *Repository) GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var org Organization
	err = conn.QueryRow(ctx,
		`SELECT id, name, description, industry, website, logo_url, timezone, currency, created_at
		 FROM org.organizations WHERE id = $1`, id).Scan(
		&org.ID, &org.Name, &org.Description, &org.Industry, &org.Website, &org.LogoURL, &org.Timezone, &org.Currency, &org.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &org, nil
}

// UpdateOrganization updates an existing organization.
func (r *Repository) UpdateOrganization(ctx context.Context, id uuid.UUID, req UpdateOrganizationRequest) (*Organization, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var org Organization
	err = conn.QueryRow(ctx,
		`UPDATE org.organizations SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			industry = COALESCE($3, industry),
			website = COALESCE($4, website),
			logo_url = COALESCE($5, logo_url),
			timezone = COALESCE($6, timezone),
			currency = COALESCE($7, currency)
			WHERE id = $8 RETURNING id, name, description, industry, website, logo_url, timezone, currency, created_at`,
		req.Name, req.Description, req.Industry, req.Website, req.LogoURL, req.Timezone, req.Currency, id).Scan(
		&org.ID, &org.Name, &org.Description, &org.Industry, &org.Website, &org.LogoURL, &org.Timezone, &org.Currency, &org.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &org, nil
}

// DeleteOrganization deletes an organization.
func (r *Repository) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM org.organizations WHERE id = $1`, id)
	return err
}

// CreateBranch creates a new branch.
func (r *Repository) CreateBranch(ctx context.Context, id uuid.UUID, req CreateBranchRequest) (*Branch, error) {
	conn, release, err := r.tenantTx(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var branch Branch
	err = conn.QueryRow(ctx,
		`INSERT INTO org.branches (id, name, description, latitude, longitude, radius_meters, ssid_list, bssid_list, tenant_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, name, description, latitude, longitude, radius_meters, ssid_list, bssid_list, tenant_id, created_at`,
		id, req.Name, req.Description, req.Latitude, req.Longitude, req.RadiusMeters, req.SSIDList, req.BSSIDList, req.TenantID, time.Now()).Scan(
		&branch.ID, &branch.Name, &branch.Description, &branch.Latitude, &branch.Longitude, &branch.RadiusMeters, 
		&branch.SSIDList, &branch.BSSIDList, &branch.TenantID, &branch.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &branch, nil
}

// ListBranches lists branches for a tenant.
func (r *Repository) ListBranches(ctx context.Context, tenantID string, page, limit int) ([]*Branch, int, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	offset := (page - 1) * limit
	
	rows, err := conn.Query(ctx,
		`SELECT id, name, description, latitude, longitude, radius_meters, ssid_list, bssid_list, tenant_id, created_at
		 FROM org.branches WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var branches []*Branch
	for rows.Next() {
		var branch Branch
		err = rows.Scan(&branch.ID, &branch.Name, &branch.Description, &branch.Latitude, &branch.Longitude, 
			&branch.RadiusMeters, &branch.SSIDList, &branch.BSSIDList, &branch.TenantID, &branch.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		branches = append(branches, &branch)
	}
	
	// Get total count
	var total int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM org.branches WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	return branches, total, nil
}

// GetBranch retrieves a branch by ID.
func (r *Repository) GetBranch(ctx context.Context, id uuid.UUID, tenantID string) (*Branch, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var branch Branch
	err = conn.QueryRow(ctx,
		`SELECT id, name, description, latitude, longitude, radius_meters, ssid_list, bssid_list, tenant_id, created_at
		 FROM org.branches WHERE id = $1 AND tenant_id = $2`, id, tenantID).Scan(
		&branch.ID, &branch.Name, &branch.Description, &branch.Latitude, &branch.Longitude, 
		&branch.RadiusMeters, &branch.SSIDList, &branch.BSSIDList, &branch.TenantID, &branch.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &branch, nil
}

// UpdateBranch updates an existing branch.
func (r *Repository) UpdateBranch(ctx context.Context, id uuid.UUID, req UpdateBranchRequest) (*Branch, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var branch Branch
	err = conn.QueryRow(ctx,
		`UPDATE org.branches SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			latitude = COALESCE($3, latitude),
			longitude = COALESCE($4, longitude),
			radius_meters = COALESCE($5, radius_meters),
			ssid_list = COALESCE($6, ssid_list),
			bssid_list = COALESCE($7, bssid_list)
			WHERE id = $8 RETURNING id, name, description, latitude, longitude, radius_meters, ssid_list, bssid_list, tenant_id, created_at`,
		req.Name, req.Description, req.Latitude, req.Longitude, req.RadiusMeters, req.SSIDList, req.BSSIDList, id).Scan(
		&branch.ID, &branch.Name, &branch.Description, &branch.Latitude, &branch.Longitude, 
		&branch.RadiusMeters, &branch.SSIDList, &branch.BSSIDList, &branch.TenantID, &branch.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &branch, nil
}

// DeleteBranch deletes a branch.
func (r *Repository) DeleteBranch(ctx context.Context, id uuid.UUID, tenantID string) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM org.branches WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// CreateEmployee creates a new employee.
func (r *Repository) CreateEmployee(ctx context.Context, id uuid.UUID, req CreateEmployeeRequest) (*Employee, error) {
	conn, release, err := r.tenantTx(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var emp Employee
	err = conn.QueryRow(ctx,
		`INSERT INTO org.employees (id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at`,
		id, req.FirstName, req.LastName, req.Email, req.EmployeeCode, req.Designation, req.Department, req.ManagerID, req.BranchID, req.TenantID, time.Now()).Scan(
		&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.EmployeeCode, &emp.Designation, &emp.Department, 
		&emp.ManagerID, &emp.BranchID, &emp.TenantID, &emp.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &emp, nil
}

// ListEmployees lists employees for a tenant.
func (r *Repository) ListEmployees(ctx context.Context, tenantID string, page, limit int) ([]*Employee, int, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	defer release()

	offset := (page - 1) * limit
	
	rows, err := conn.Query(ctx,
		`SELECT id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at
		 FROM org.employees WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var employees []*Employee
	for rows.Next() {
		var emp Employee
		err = rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.EmployeeCode, &emp.Designation, 
			&emp.Department, &emp.ManagerID, &emp.BranchID, &emp.TenantID, &emp.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		employees = append(employees, &emp)
	}
	
	// Get total count
	var total int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM org.employees WHERE tenant_id = $1`, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	
	return employees, total, nil
}

// GetEmployee retrieves an employee by ID.
func (r *Repository) GetEmployee(ctx context.Context, id uuid.UUID, tenantID string) (*Employee, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var emp Employee
	err = conn.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at
		 FROM org.employees WHERE id = $1 AND tenant_id = $2`, id, tenantID).Scan(
		&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.EmployeeCode, &emp.Designation, 
		&emp.Department, &emp.ManagerID, &emp.BranchID, &emp.TenantID, &emp.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &emp, nil
}

// UpdateEmployee updates an existing employee.
func (r *Repository) UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*Employee, error) {
	conn, release, err := r.tenantTx(ctx, "")
	if err != nil {
		return nil, err
	}
	defer release()

	var emp Employee
	err = conn.QueryRow(ctx,
		`UPDATE org.employees SET first_name = COALESCE($1, first_name),
			last_name = COALESCE($2, last_name),
			email = COALESCE($3, email),
			employee_code = COALESCE($4, employee_code),
			designation = COALESCE($5, designation),
			department = COALESCE($6, department),
			manager_id = COALESCE($7, manager_id),
			branch_id = COALESCE($8, branch_id)
			WHERE id = $9 RETURNING id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at`,
		req.FirstName, req.LastName, req.Email, req.EmployeeCode, req.Designation, req.Department, req.ManagerID, req.BranchID, id).Scan(
		&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.EmployeeCode, &emp.Designation, 
		&emp.Department, &emp.ManagerID, &emp.BranchID, &emp.TenantID, &emp.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &emp, nil
}

// DeleteEmployee deletes an employee.
func (r *Repository) DeleteEmployee(ctx context.Context, id uuid.UUID, tenantID string) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`DELETE FROM org.employees WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// GetEmployeeByCode retrieves an employee by their employee code.
func (r *Repository) GetEmployeeByCode(ctx context.Context, tenantID, employeeCode string) (*Employee, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var emp Employee
	err = conn.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at
		 FROM org.employees WHERE employee_code = $1 AND tenant_id = $2`, employeeCode, tenantID).Scan(
		&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.EmployeeCode, &emp.Designation, 
		&emp.Department, &emp.ManagerID, &emp.BranchID, &emp.TenantID, &emp.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &emp, nil
}

// CreateEmployeeWithDocuments creates an employee with KYC documents.
func (r *Repository) CreateEmployeeWithDocuments(ctx context.Context, id uuid.UUID, req CreateEmployeeWithDocumentsRequest) (*EmployeeWithDocuments, error) {
	conn, release, err := r.tenantTx(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var emp EmployeeWithDocuments
	err = conn.QueryRow(ctx,
		`INSERT INTO org.employees (id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, first_name, last_name, email, employee_code, designation, department, manager_id, branch_id, tenant_id, created_at`,
		id, req.FirstName, req.LastName, req.Email, req.EmployeeCode, req.Designation, req.Department, req.ManagerID, req.BranchID, req.TenantID, time.Now()).Scan(
		&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.EmployeeCode, &emp.Designation, &emp.Department, 
		&emp.ManagerID, &emp.BranchID, &emp.TenantID, &emp.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	// In a real implementation, we would also save the documents here
	
	return &emp, nil
}

// GetEmployeeDocuments retrieves employee KYC documents.
func (r *Repository) GetEmployeeDocuments(ctx context.Context, employeeID uuid.UUID, tenantID string) (*EmployeeDocuments, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	// In a real implementation, this would query documents table
	var docs EmployeeDocuments
	
	// For now, return mock data
	docs = EmployeeDocuments{
		EmployeeID: employeeID,
		PanNumber:  "PAN123456789",
		AadhaarNumber: "123456789012",
		UanNumber: "UAN9876543210",
		EsicNumber: "ESIC123456789",
		Documents: []Document{
			{
				Name: "PAN Card",
				URL:  "https://example.com/documents/pan_card.pdf",
			},
			{
				Name: "Aadhaar Card",
				URL:  "https://example.com/documents/aadhaar_card.pdf",
			},
		},
	}
	
	return &docs, nil
}

// UpdateEmployeeExitStatus updates employee status to exited.
func (r *Repository) UpdateEmployeeExitStatus(ctx context.Context, employeeID uuid.UUID, tenantID string, exitStatus EmployeeExitStatus) error {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer release()

	_, err = conn.Exec(ctx,
		`UPDATE org.employees SET status = $1, last_working_day = $2, exit_reason = $3
		 WHERE id = $4 AND tenant_id = $5`,
		exitStatus.Status, exitStatus.LastWorkingDay, exitStatus.ExitReason, employeeID, tenantID)
	return err
}

// CreateSalaryStructure creates a salary structure for an employee.
func (r *Repository) CreateSalaryStructure(ctx context.Context, employeeID uuid.UUID, tenantID string, salaryReq CreateSalaryStructureRequest) (*SalaryStructure, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var salaryStruct SalaryStructure
	err = conn.QueryRow(ctx,
		`INSERT INTO org.salary_structures (employee_id, tenant_id, pay_basis, basic, hra, da, conveyance, special_allowance, 
			medical_allowance, bonus, overtime_rate, effective_date, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING employee_id, tenant_id, pay_basis, basic, hra, da, conveyance, 
			special_allowance, medical_allowance, bonus, overtime_rate, effective_date, created_at`,
		employeeID, tenantID, salaryReq.PayBasis, salaryReq.Basic, salaryReq.HRA, salaryReq.DA, salaryReq.Conveyance,
		salaryReq.SpecialAllowance, salaryReq.MedicalAllowance, salaryReq.Bonus, salaryReq.OvertimeRate, salaryReq.EffectiveDate, time.Now()).Scan(
		&salaryStruct.EmployeeID, &salaryStruct.TenantID, &salaryStruct.PayBasis, &salaryStruct.Basic, &salaryStruct.HRA, 
		&salaryStruct.DA, &salaryStruct.Conveyance, &salaryStruct.SpecialAllowance, &salaryStruct.MedicalAllowance, 
		&salaryStruct.Bonus, &salaryStruct.OvertimeRate, &salaryStruct.EffectiveDate, &salaryStruct.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &salaryStruct, nil
}

// GetSalaryStructure retrieves an employee's salary structure.
func (r *Repository) GetSalaryStructure(ctx context.Context, employeeID uuid.UUID, tenantID string) (*SalaryStructure, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var salaryStruct SalaryStructure
	err = conn.QueryRow(ctx,
		`SELECT employee_id, tenant_id, pay_basis, basic, hra, da, conveyance, 
			special_allowance, medical_allowance, bonus, overtime_rate, effective_date, created_at
		 FROM org.salary_structures WHERE employee_id = $1 AND tenant_id = $2`, employeeID, tenantID).Scan(
		&salaryStruct.EmployeeID, &salaryStruct.TenantID, &salaryStruct.PayBasis, &salaryStruct.Basic, &salaryStruct.HRA, 
		&salaryStruct.DA, &salaryStruct.Conveyance, &salaryStruct.SpecialAllowance, &salaryStruct.MedicalAllowance, 
		&salaryStruct.Bonus, &salaryStruct.OvertimeRate, &salaryStruct.EffectiveDate, &salaryStruct.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	return &salaryStruct, nil
}

// UpdateSalaryStructure updates an employee's salary structure.
func (r *Repository) UpdateSalaryStructure(ctx context.Context, employeeID uuid.UUID, tenantID string, salaryReq CreateSalaryStructureRequest) (*SalaryStructure, error) {
	conn, release, err := r.tenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer release()

	var salaryStruct SalaryStructure
	err = conn.QueryRow(ctx,
		`UPDATE org.salary_structures SET pay_basis = COALESCE($1, pay_basis),
			basic = COALESCE($2, basic),
			hra = COALESCE($3, hra),
			da = COALESCE($4, da),
			conveyance = COALESCE($5, conveyance),
			special_allowance = COALESCE($6, special_allowance),
			medical_allowance = COALESCE($7, medical_allowance),
			bonus = COALESCE($8, bonus),
			overtime_rate = COALESCE($9, overtime_rate),
			effective_date = COALESCE($10, effective_date),
			updated_at = NOW()
			WHERE employee_id = $11 AND tenant_id = $12 RETURNING employee_id, tenant_id, pay_basis, basic, hra, da, conveyance, 
			special_allowance, medical_allowance, bonus, overtime_rate, effective_date, created_at, updated_at`,
		salaryReq.PayBasis, salaryReq.Basic, salaryReq.HRA, salaryReq.DA, salaryReq.Conveyance,
		salaryReq.SpecialAllowance, salaryReq.MedicalAllowance, salaryReq.Bonus, salaryReq.OvertimeRate, salaryReq.EffectiveDate,
		employeeID, tenantID).Scan(
		&salaryStruct.EmployeeID, &salaryStruct.TenantID, &salaryStruct.PayBasis, &salaryStruct.Basic, &salaryStruct.HRA, 
		&salaryStruct.DA, &salaryStruct.Conveyance, &salaryStruct.SpecialAllowance, &salaryStruct.MedicalAllowance, 
		&salaryStruct.Bonus, &salaryStruct.OvertimeRate, &salaryStruct.EffectiveDate, &salaryStruct.CreatedAt, &salaryStruct.UpdatedAt)
	if err != nil {
		return nil, err
	}
	
	return &salaryStruct, nil
}
