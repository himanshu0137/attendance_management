package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service handles business logic for users
type Service struct {
	repo *Repository
}

// NewService creates a new user service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetMe gets the current user's details
func (s *Service) GetMe(ctx context.Context, tenantID, userID string) (*User, error) {
	user, err := s.repo.GetUser(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	
	// If user doesn't exist in our database, create a default one
	if user == nil {
		// This would typically be handled by a registration process
		// For now we return a default user or error
		return nil, nil
	}
	
	return user, nil
}

// UpdateMe updates the current user's details
func (s *Service) UpdateMe(ctx context.Context, tenantID, userID string, req UpdateUserReq) (*User, error) {
	user := &User{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Language: req.Language,
	}
	
	err := s.repo.UpdateUser(ctx, tenantID, userID, user)
	if err != nil {
		return nil, err
	}
	
	return s.GetMe(ctx, tenantID, userID)
}

// GetPreferences gets the current user's preferences
func (s *Service) GetPreferences(ctx context.Context, tenantID, userID string) (*Preferences, error) {
	prefs, err := s.repo.GetUserPreferences(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	
	return prefs, nil
}

// UpdatePreferences updates the current user's preferences
func (s *Service) UpdatePreferences(ctx context.Context, tenantID, userID string, prefs Preferences) (*Preferences, error) {
	err := s.repo.UpdateUserPreferences(ctx, tenantID, userID, &prefs)
	if err != nil {
		return nil, err
	}
	
	return s.GetPreferences(ctx, tenantID, userID)
}

// InitializeUser creates a default user record for a new user
func (s *Service) InitializeUser(ctx context.Context, tenantID, keycloakSub string, name, email, role string) (*User, error) {
	user := &User{
		ID:          uuid.New().String(),
		KeycloakSub: keycloakSub,
		Name:        name,
		Email:       email,
		Role:        UserRole(role), // Convert to our role type
		Language:    "en", // Default language
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	
	err := s.repo.CreateUser(ctx, tenantID, keycloakSub, user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}