package leave

import (
	"context"
	"database/sql"
	"time"

	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateLeaveRequest(ctx context.Context, leave LeaveRequest) error
	GetLeaveRequest(ctx context.Context, id uuid.UUID) (*LeaveRequest, error)
	ListLeaveRequests(ctx context.Context, tenantID, status string, page, limit int) ([]*LeaveRequest, int, error)
	UpdateLeaveStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *string, approvedAt *time.Time) error
	ListMyLeaveRequests(ctx context.Context, employeeID string, page, limit int) ([]*LeaveRequest, int, error)
	// TODO: Add methods for leave balance checking and management
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateLeaveRequest(ctx context.Context, leave LeaveRequest) error {
	query := `
		INSERT INTO leave_requests (id, employee_id, tenant_id, start_date, end_date, leave_type, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query, leave.ID, leave.EmployeeID, leave.TenantID, leave.StartDate, leave.EndDate, leave.LeaveType, leave.Description, leave.Status, leave.CreatedAt, leave.UpdatedAt)
	return err
}

func (r *repository) GetLeaveRequest(ctx context.Context, id uuid.UUID) (*LeaveRequest, error) {
	query := `
		SELECT id, employee_id, tenant_id, start_date, end_date, leave_type, description, status, approved_by, approved_at, created_at, updated_at
		FROM leave_requests
		WHERE id = $1
	`

	var leave LeaveRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&leave.ID,
		&leave.EmployeeID,
		&leave.TenantID,
		&leave.StartDate,
		&leave.EndDate,
		&leave.LeaveType,
		&leave.Description,
		&leave.Status,
		&leave.ApprovedBy,
		&leave.ApprovedAt,
		&leave.CreatedAt,
		&leave.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &leave, nil
}

func (r *repository) ListLeaveRequests(ctx context.Context, tenantID, status string, page, limit int) ([]*LeaveRequest, int, error) {
	offset := (page - 1) * limit

	var query string
	var args []interface{}

	if status != "" {
		query = `
			SELECT id, employee_id, tenant_id, start_date, end_date, leave_type, description, status, approved_by, approved_at, created_at, updated_at
			FROM leave_requests
			WHERE tenant_id = $1 AND status = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{tenantID, status, limit, offset}
	} else {
		query = `
			SELECT id, employee_id, tenant_id, start_date, end_date, leave_type, description, status, approved_by, approved_at, created_at, updated_at
			FROM leave_requests
			WHERE tenant_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{tenantID, limit, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leaves []*LeaveRequest
	for rows.Next() {
		var leave LeaveRequest
		err := rows.Scan(
			&leave.ID,
			&leave.EmployeeID,
			&leave.TenantID,
			&leave.StartDate,
			&leave.EndDate,
			&leave.LeaveType,
			&leave.Description,
			&leave.Status,
			&leave.ApprovedBy,
			&leave.ApprovedAt,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		leaves = append(leaves, &leave)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM leave_requests
		WHERE tenant_id = $1
	`
	var total int
	err = r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

func (r *repository) UpdateLeaveStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *string, approvedAt *time.Time) error {
	query := `
		UPDATE leave_requests
		SET status = $1, approved_by = $2, approved_at = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query, status, approvedBy, approvedAt, time.Now(), id)
	return err
}

func (r *repository) ListMyLeaveRequests(ctx context.Context, employeeID string, page, limit int) ([]*LeaveRequest, int, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, employee_id, tenant_id, start_date, end_date, leave_type, description, status, approved_by, approved_at, created_at, updated_at
		FROM leave_requests
		WHERE employee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, employeeID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leaves []*LeaveRequest
	for rows.Next() {
		var leave LeaveRequest
		err := rows.Scan(
			&leave.ID,
			&leave.EmployeeID,
			&leave.TenantID,
			&leave.StartDate,
			&leave.EndDate,
			&leave.LeaveType,
			&leave.Description,
			&leave.Status,
			&leave.ApprovedBy,
			&leave.ApprovedAt,
			&leave.CreatedAt,
			&leave.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		leaves = append(leaves, &leave)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM leave_requests
		WHERE employee_id = $1
	`
	var total int
	err = r.db.QueryRow(ctx, countQuery, employeeID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}
