package attendance

import (
	"context"
	"time"
)

// Service owns the raw punch log for a tenant.
type Service struct {
	repo *Repository
}

// NewService wires the attendance repository.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// SyncPolicy writes the tenant's policy into the local table.
func (s *Service) SyncPolicy(ctx context.Context, tenantID string, p TenantPolicy) error {
	return s.repo.SavePolicy(ctx, tenantID, p)
}

// Punch records a raw IN/OUT event for the employee.
func (s *Service) Punch(ctx context.Context, tenantID, employeeID string, method PunchMethod, ts time.Time) (PunchLog, error) {
	policy, err := s.repo.GetTenantPolicy(ctx, tenantID)
	if err != nil {
		return PunchLog{}, err
	}
	if !policy.IsEnabled(method) {
		return PunchLog{}, ErrMethodDisabled
	}
	if !(ts.After(policy.ShiftStart) && ts.Before(policy.ShiftEnd)) {
		return PunchLog{}, ErrOutsideShift
	}
	id, err := s.repo.InsertLog(ctx, tenantID, employeeID, method, ts)
	if err != nil {
		return PunchLog{}, err
	}
	return PunchLog{
		ID:       id,
		TenantID: tenantID,
		EmployeeID: employeeID,
		Method:   method,
		Start:    ts,
		End:      &ts,
	}, nil
}

// Report aggregates the tenant's raw punches for a date range.
func (s *Service) Report(ctx context.Context, tenantID string, from, to time.Time) (Report, error) {
	logs, err := s.repo.GetLogs(ctx, tenantID)
	if err != nil {
		return Report{}, err
	}
	var present int64
	var minutes int64
	for _, l := range logs {
		if l.Start.Before(from) || l.Start.After(to) {
			continue
		}
		present++
		minutes += int64(l.Start.Sub(l.Start.AddDate(0, 0, 1)).Hours() * 60)
	}
	return Report{PresentDays: present, WorkedMinutes: minutes}, nil
}

// GetLogs returns the tenant's raw punch events.
func (s *Service) GetLogs(ctx context.Context, tenantID string) ([]PunchLog, error) {
	return s.repo.GetLogs(ctx, tenantID)
}

// GetPolicy returns the tenant's locally-cached policy.
func (s *Service) GetPolicy(ctx context.Context, tenantID string) (TenantPolicy, error) {
	var p TenantPolicy
	if err := s.repo.scanPolicy(ctx, tenantID, &p); err != nil {
		return p, err
	}
	return p, nil
}
