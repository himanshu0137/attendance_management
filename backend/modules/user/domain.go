package user

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

// User represents a system user with their details and role
type User struct {
	ID          string     `json:"id"`
	KeycloakSub string     `json:"keycloak_sub"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Role        UserRole   `json:"role"`
	Language    string     `json:"language"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
}

// Preferences represents user application preferences
type Preferences struct {
	Theme         string `json:"theme"`           // light, dark
	Notifications bool   `json:"notifications"`   // enable/disable notifications
	Language      string `json:"language"`        // UI language
	DateFormat    string `json:"date_format"`     // date format preference
	AutoClockIn   bool   `json:"auto_clock_in"`   // auto clock in on app launch
}