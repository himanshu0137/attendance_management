package constants

// UserRole represents the different roles a user can have in the system
type UserRole string

const (
	// UserRoleAdmin represents an administrator role
	UserRoleAdmin UserRole = "admin"

	// UserRoleManager represents a manager/HR role
	UserRoleManager UserRole = "manager"

	// UserRoleEmployee represents a regular employee role
	UserRoleEmployee UserRole = "employee"

	// UserRoleGuest represents a guest user role
	UserRoleGuest UserRole = "guest"
)
