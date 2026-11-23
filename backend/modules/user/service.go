package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// User represents a user in the system
type User struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UserCreateRequest represents a request to create a new user
type UserCreateRequest struct {
	TenantID      uuid.UUID `json:"tenant_id"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
}

// UserUpdateRequest represents a request to update a user
type UserUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

// Role represents a role in the system
type Role struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRoleAssignment represents the assignment of a role to a user
type UserRoleAssignment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	RoleID      uuid.UUID `json:"role_id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	AssignedAt  time.Time `json:"assigned_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateUserService creates a new user
func (s *Service) CreateUser(ctx context.Context, req UserCreateRequest) (*User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := User{
		ID:           uuid.New(),
		TenantID:     req.TenantID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// GetUser gets a user by ID
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user
func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req UserUpdateRequest) error {
	err := s.repo.UpdateUser(ctx, id, req)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers lists users with pagination
func (s *Service) ListUsers(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*User, int, error) {
	users, total, err := s.repo.ListUsers(ctx, tenantID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	return users, total, nil
}

// CreateRole creates a new role
func (s *Service) CreateRole(ctx context.Context, tenantID uuid.UUID, name, description string, permissions []string) (*Role, error) {
	role := Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Permissions: permissions,
		CreatedAt:   time.Now(),
	}

	err := s.repo.CreateRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

// GetRole gets a role by ID
func (s *Service) GetRole(ctx context.Context, id uuid.UUID) (*Role, error) {
	role, err := s.repo.GetRole(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return role, nil
}

// AssignRole assigns a role to a user
func (s *Service) AssignRole(ctx context.Context, userID, roleID, tenantID uuid.UUID) error {
	assignment := UserRoleAssignment{
		ID:         uuid.New(),
		UserID:     userID,
		RoleID:     roleID,
		TenantID:   tenantID,
		AssignedAt: time.Now(),
		CreatedAt:  time.Now(),
	}

	err := s.repo.AssignRole(ctx, assignment)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}
	return nil
}

// GetUserRoles gets all roles assigned to a user
func (s *Service) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*UserRoleAssignment, error) {
	assignments, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	return assignments, nil
}

// CheckPermission checks if a user has a specific permission
func (s *Service) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	hasPermission, err := s.repo.CheckPermission(ctx, userID, permission)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return hasPermission, nil
}