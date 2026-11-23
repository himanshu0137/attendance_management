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
		ID:         id,
		TenantID:   tenantID,
		EmployeeID: employeeID,
		Method:     method,
		Start:      ts,
		End:        &ts,
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

// CreateShift creates a new shift for a tenant
func (s *Service) CreateShift(ctx context.Context, tenantID string, shift Shift) (*Shift, error) {
	createdShift, err := s.repo.CreateShift(ctx, tenantID, shift)
	if err != nil {
		return nil, err
	}
	return createdShift, nil
}

// GetShift retrieves a shift by ID
func (s *Service) GetShift(ctx context.Context, tenantID string, shiftID string) (*Shift, error) {
	shift, err := s.repo.GetShift(ctx, tenantID, shiftID)
	if err != nil {
		return nil, err
	}
	return shift, nil
}

// ListShifts lists shifts for a tenant
func (s *Service) ListShifts(ctx context.Context, tenantID string, page, limit int) ([]*Shift, int, error) {
	shifts, total, err := s.repo.ListShifts(ctx, tenantID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return shifts, total, nil
}

// UpdateShift updates an existing shift
func (s *Service) UpdateShift(ctx context.Context, tenantID string, shiftID string, shift Shift) (*Shift, error) {
	updatedShift, err := s.repo.UpdateShift(ctx, tenantID, shiftID, shift)
	if err != nil {
		return nil, err
	}
	return updatedShift, nil
}

// DeleteShift deletes a shift
func (s *Service) DeleteShift(ctx context.Context, tenantID string, shiftID string) error {
	return s.repo.DeleteShift(ctx, tenantID, shiftID)
}

// CreateAttendanceRule creates a new attendance rule for a tenant
func (s *Service) CreateAttendanceRule(ctx context.Context, tenantID string, rule AttendanceRule) (*AttendanceRule, error) {
	createdRule, err := s.repo.CreateAttendanceRule(ctx, tenantID, rule)
	if err != nil {
		return nil, err
	}
	return createdRule, nil
}

// GetAttendanceRule retrieves an attendance rule by ID
func (s *Service) GetAttendanceRule(ctx context.Context, tenantID string, ruleID string) (*AttendanceRule, error) {
	rule, err := s.repo.GetAttendanceRule(ctx, tenantID, ruleID)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// ListAttendanceRules lists attendance rules for a tenant
func (s *Service) ListAttendanceRules(ctx context.Context, tenantID string, page, limit int) ([]*AttendanceRule, int, error) {
	rules, total, err := s.repo.ListAttendanceRules(ctx, tenantID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// UpdateAttendanceRule updates an existing attendance rule
func (s *Service) UpdateAttendanceRule(ctx context.Context, tenantID string, ruleID string, rule AttendanceRule) (*AttendanceRule, error) {
	updatedRule, err := s.repo.UpdateAttendanceRule(ctx, tenantID, ruleID, rule)
	if err != nil {
		return nil, err
	}
	return updatedRule, nil
}

// DeleteAttendanceRule deletes an attendance rule
func (s *Service) DeleteAttendanceRule(ctx context.Context, tenantID string, ruleID string) error {
	return s.repo.DeleteAttendanceRule(ctx, tenantID, ruleID)
}
