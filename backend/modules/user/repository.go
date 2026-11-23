package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for users
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new user repository
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetUser gets a user by tenant ID and user ID
func (r *Repository) GetUser(ctx context.Context, tenantID, userID string) (*User, error) {
	query := `
		SELECT id, keycloak_sub, name, email, phone, role, language, created_at, updated_at
		FROM user.users 
		WHERE tenant_id = $1 AND id = $2
	`
	
	var user User
	err := r.pool.QueryRow(ctx, query, tenantID, userID).Scan(
		&user.ID,
		&user.KeycloakSub,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.Language,
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

// UpdateUser updates a user's details
func (r *Repository) UpdateUser(ctx context.Context, tenantID, userID string, user *User) error {
	query := `
		UPDATE user.users 
		SET name = $1, email = $2, phone = $3, language = $4, updated_at = NOW()
		WHERE tenant_id = $5 AND id = $6
	`
	
	_, err := r.pool.Exec(ctx, query, user.Name, user.Email, user.Phone, user.Language, tenantID, userID)
	return err
}

// GetUserPreferences gets a user's preferences
func (r *Repository) GetUserPreferences(ctx context.Context, tenantID, userID string) (*Preferences, error) {
	query := `
		SELECT theme, notifications, language, date_format, auto_clock_in
		FROM user.preferences 
		WHERE tenant_id = $1 AND user_id = $2
	`
	
	var prefs Preferences
	err := r.pool.QueryRow(ctx, query, tenantID, userID).Scan(
		&prefs.Theme,
		&prefs.Notifications,
		&prefs.Language,
		&prefs.DateFormat,
		&prefs.AutoClockIn,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return &Preferences{}, nil
		}
		return nil, err
	}
	
	return &prefs, nil
}

// UpdateUserPreferences updates a user's preferences
func (r *Repository) UpdateUserPreferences(ctx context.Context, tenantID, userID string, prefs *Preferences) error {
	query := `
		INSERT INTO user.preferences (tenant_id, user_id, theme, notifications, language, date_format, auto_clock_in)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tenant_id, user_id) 
		DO UPDATE SET 
			theme = EXCLUDED.theme,
			notifications = EXCLUDED.notifications,
			language = EXCLUDED.language,
			date_format = EXCLUDED.date_format,
			auto_clock_in = EXCLUDED.auto_clock_in,
			updated_at = NOW()
	`
	
	_, err := r.pool.Exec(ctx, query, tenantID, userID, prefs.Theme, prefs.Notifications, prefs.Language, prefs.DateFormat, prefs.AutoClockIn)
	return err
}

// CreateUser creates a new user record
func (r *Repository) CreateUser(ctx context.Context, tenantID, keycloakSub string, user *User) error {
	query := `
		INSERT INTO user.users (tenant_id, id, keycloak_sub, name, email, phone, role, language, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	
	_, err := r.pool.Exec(ctx, query, tenantID, user.ID, keycloakSub, user.Name, user.Email, user.Phone, user.Role, user.Language)
	return err
}