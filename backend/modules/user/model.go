package user

import (
	"uuid"

	"hunrx.com/attendance_management/pkg/constants"
)

type User struct {
	ID         uuid.UUID          `db:"id"`
	TenantID   uuid.UUID          `db:"tenant_id"`
	KeyCloakID string             `db:"keycloak_id"`
	Name       string             `db:"name"`
	Email      string             `db:"email"`
	Phone      string             `db:"phone"`
	Language   string             `db:"language"`
	Role       constants.UserRole `db:"role"`
	CreatedAt  string             `db:"created_at"`
	UpdatedAt  string             `db:"updated_at"`
}
