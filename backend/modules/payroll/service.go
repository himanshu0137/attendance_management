package payroll

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// Service owns the payroll schema and conditional slip generation.
type Service struct {
	repo *Repository
}

// NewService wires the payroll repository.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetSettings returns the tenant's salary slip eligibility.
func (s *Service) GetSettings(ctx context.Context, tenantID string) (TenantSettings, error) {
	return s.repo.GetTenantSettings(ctx, tenantID)
}

// SaveSettings writes the tenant's salary slip eligibility.
func (s *Service) SaveSettings(ctx context.Context, tenantID string, t TenantSettings) error {
	return s.repo.SaveSettings(ctx, tenantID, t)
}

// IsEligible reports whether payroll may generate slips for the tenant.
func (s *Service) IsEligible(ctx context.Context, tenantID string) (bool, error) {
	t, err := s.repo.GetTenantSettings(ctx, tenantID)
	if err != nil {
		return false, err
	}
	return t.SalarySlipEnabled, nil
}

// GenerateSlip computes a prorated monthly salary slip for the employee.
func (s *Service) GenerateSlip(ctx context.Context, tenantID, employeeID string, month time.Month, year int) (SalarySlip, error) {
	settings, err := s.repo.GetTenantSettings(ctx, tenantID)
	if err != nil {
		return SalarySlip{}, err
	}
	if !settings.SalarySlipEnabled {
		return SalarySlip{}, ErrNotEligible
	}

	structures, err := s.repo.GetStructures(ctx, tenantID)
	if err != nil {
		return SalarySlip{}, err
	}
	var structure SalaryStructure
	for _, st := range structures {
		if st.EmployeeID == employeeID {
			structure = st
		}
	}
	if structure.EmployeeID == "" {
		return SalarySlip{}, ErrNoStructure
	}

	days, base := s.prorate(structure, month, year)
	period := fmt.Sprintf("%d-%s", year, month.String())
	id := fmt.Sprintf("slip-%s-%d-%s", employeeID, year, month.String())
	return SalarySlip{
		ID:         id,
		TenantID:   tenantID,
		EmployeeID: employeeID,
		Period:     period,
		BaseSalary: base,
		Days:       days,
	}, nil
}

// InsertSlip writes a generated salary slip.
func (s *Service) InsertSlip(ctx context.Context, tenantID, employeeID string, period string, slip SalarySlip) error {
	return s.repo.InsertSlip(ctx, tenantID, employeeID, period, slip)
}

// GetSlips returns the tenant's generated salary slips.
func (s *Service) GetSlips(ctx context.Context, tenantID string) ([]SalarySlip, error) {
	return s.repo.GetSlips(ctx, tenantID)
}

// prorate splits the monthly salary across the payable days of the month.
func (s *Service) prorate(structure SalaryStructure, month time.Month, year int) ([]Day, Money) {
	daysPerMonth := int64(structure.DaysPerMonth)
	if daysPerMonth == 0 {
		daysPerMonth = int64(daysInMonth(month, year))
	}
	total := NewMoney(structure.BaseSalary.String())
	if total.Cmp(NewMoney("0")) == 0 {
		return make([]Day, 0, daysPerMonth), Money{}
	}
	perDay := total.Div(NewMoney(strconv.FormatInt(daysPerMonth, 10)))
	days := make([]Day, 0, daysPerMonth)
	paid := Money{}
	for d := 1; d <= int(daysPerMonth); d++ {
		days = append(days, Day{Date: time.Date(year, month, d, 0, 0, 0, 0, time.UTC)})
		if total.Sub(paid).Cmp(NewMoney("1")) <= 0 {
			break
		}
		paid = paid.Add(perDay)
	}
	return days, paid
}

// daysInMonth returns the number of days in the month.
func daysInMonth(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
