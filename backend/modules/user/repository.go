package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UserUpdateRequest) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*User, int, error)
	
	CreateRole(ctx context.Context, role Role) error
	GetRole(ctx context.Context, id uuid.UUID) (*Role, error)
	AssignRole(ctx context.Context, assignment UserRoleAssignment) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*UserRoleAssignment, error)
	CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(ctx, query, user.ID, user.TenantID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.TenantID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *repository) UpdateUser(ctx context.Context, id uuid.UUID, req UserUpdateRequest) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(ctx, query, req.FirstName, req.LastName, req.IsActive, time.Now(), id)
	return err
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListUsers(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*User, int, error) {
	offset := (page - 1) * limit
	
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.TenantID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE tenant_id = $1
	`
	var total int
	err = r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) CreateRole(ctx context.Context, role Role) error {
	query := `
		INSERT INTO roles (id, tenant_id, name, description, permissions, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, role.ID, role.TenantID, role.Name, role.Description, role.Permissions, role.CreatedAt)
	return err
}

func (r *repository) GetRole(ctx context.Context, id uuid.UUID) (*Role, error) {
	query := `
		SELECT id, tenant_id, name, description, permissions, created_at
		FROM roles
		WHERE id = $1
	`

	var role Role
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.TenantID,
		&role.Name,
		&role.Description,
		&role.Permissions,
		&role.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &role, nil
}

func (r *repository) AssignRole(ctx context.Context, assignment UserRoleAssignment) error {
	query := `
		INSERT INTO user_role_assignments (id, user_id, role_id, tenant_id, assigned_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, assignment.ID, assignment.UserID, assignment.RoleID, assignment.TenantID, assignment.AssignedAt, assignment.CreatedAt)
	return err
}

func (r *repository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*UserRoleAssignment, error) {
	query := `
		SELECT id, user_id, role_id, tenant_id, assigned_at, created_at
		FROM user_role_assignments
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*UserRoleAssignment
	for rows.Next() {
		var assignment UserRoleAssignment
		err := rows.Scan(
			&assignment.ID,
			&assignment.UserID,
			&assignment.RoleID,
			&assignment.TenantID,
			&assignment.AssignedAt,
			&assignment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}

	return assignments, nil
}

func (r *repository) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_role_assignments ura
			JOIN roles r ON ura.role_id = r.id
			WHERE ura.user_id = $1
			AND $2 = ANY(r.permissions)
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, permission).Scan(&exists)
	if err != nil {
		return false, err
	}
	
	return exists, nil
}