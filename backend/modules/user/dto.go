package user

import (
	"time"
	"uuid"
)

// Preferences represents user application preferences
type Preferences struct {
	Theme         string `json:"theme"`         // light, dark
	Notifications bool   `json:"notifications"` // enable/disable notifications
	Language      string `json:"language"`      // UI language
	DateFormat    string `json:"date_format"`   // date format preference
	AutoClockIn   bool   `json:"auto_clock_in"` // auto clock in on app launch
}

// User represents a user in the system
type User struct {
	ID uuid.UUID `json:"id"`

	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserCreateRequest represents a request to create a new user
type UserCreateRequest struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
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
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	AssignedAt time.Time `json:"assigned_at"`
	CreatedAt  time.Time `json:"created_at"`
}
