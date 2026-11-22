package organization

import (
	"context"

	"hunrx.com/attendance_management/pkg/events"
)

// Service owns the org schema write path and the policy event.
type Service struct {
	repo     *Repository
	bus      events.Bus
	keycloak *KeycloakClient
}

// KeycloakClient wraps the headless Keycloak admin REST client.
type KeycloakClient struct {
	Issuer, Realm, ClientID, ClientSecret string
}

// NewService wires the org repository, event bus and Keycloak client.
func NewService(repo *Repository, bus events.Bus, kc *KeycloakClient) *Service {
	return &Service{repo: repo, bus: bus, keycloak: kc}
}

// SetPolicy writes the tenant's policy and publishes org.policy.updated.
func (s *Service) SetPolicy(ctx context.Context, tenantID string, p OrgPolicy) error {
	if err := s.repo.SavePolicy(ctx, tenantID, p); err != nil {
		return err
	}
	s.bus.Publish(events.Event{Name: "org.policy.updated", Data: p})
	return nil
}

// GetPolicy returns the tenant's policy, or ErrPolicyNotFound when unset.
func (s *Service) GetPolicy(ctx context.Context, tenantID string) (OrgPolicy, error) {
	return s.repo.GetPolicy(ctx, tenantID)
}

// CreateEmployee saves a Keycloak user for the tenant and returns its sub.
func (s *Service) CreateEmployee(ctx context.Context, tenantID, username, role string) (string, error) {
	sub, err := CreateKeycloakUser(ctx, s.keycloak.Issuer, s.keycloak.Realm, s.keycloak.ClientID, s.keycloak.ClientSecret, username, role)
	if err != nil {
		return "", err
	}
	return s.repo.CreateEmployee(ctx, tenantID, sub, role)
}

// ListEmployees lists the tenant's Keycloak-backed users.
func (s *Service) ListEmployees(ctx context.Context, tenantID string) ([]Employee, error) {
	return s.repo.GetEmployees(ctx, tenantID)
}
