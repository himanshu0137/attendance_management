package leave

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateLeaveRequest creates a new leave request
func (s *Service) CreateLeaveRequest(ctx context.Context, req CreateLeaveRequest) (*LeaveRequest, error) {
	// TODO: Add leave balance validation here for TC-13.1
	leave := LeaveRequest{
		ID:          uuid.New(),
		EmployeeID:  req.EmployeeID,
		TenantID:    req.TenantID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		LeaveType:   req.LeaveType,
		Description: req.Description,
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.CreateLeaveRequest(ctx, leave)
	if err != nil {
		return nil, err
	}

	return &leave, nil
}

// GetLeaveRequest gets a leave request by ID
func (s *Service) GetLeaveRequest(ctx context.Context, id uuid.UUID) (*LeaveRequest, error) {
	leave, err := s.repo.GetLeaveRequest(ctx, id)
	if err != nil {
		return nil, err
	}
	return leave, nil
}

// ListLeaveRequests lists leave requests with optional filtering
func (s *Service) ListLeaveRequests(ctx context.Context, tenantID, status string, page, limit int) ([]*LeaveRequest, int, error) {
	leaves, total, err := s.repo.ListLeaveRequests(ctx, tenantID, status, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return leaves, total, nil
}

// ApproveLeaveRequest approves a leave request
func (s *Service) ApproveLeaveRequest(ctx context.Context, id uuid.UUID, approvedBy string) error {
	approvedAt := time.Now()
	
	err := s.repo.UpdateLeaveStatus(ctx, id, "APPROVED", &approvedBy, &approvedAt)
	if err != nil {
		return err
	}
	
	return nil
}

// RejectLeaveRequest rejects a leave request
func (s *Service) RejectLeaveRequest(ctx context.Context, id uuid.UUID, rejectedBy string) error {
	rejectedAt := time.Now()
	
	err := s.repo.UpdateLeaveStatus(ctx, id, "REJECTED", &rejectedBy, &rejectedAt)
	if err != nil {
		return err
	}
	
	return nil
}

// ListMyLeaveRequests lists leave requests for a specific employee
func (s *Service) ListMyLeaveRequests(ctx context.Context, employeeID string, page, limit int) ([]*LeaveRequest, int, error) {
	leaves, total, err := s.repo.ListMyLeaveRequests(ctx, employeeID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return leaves, total, nil
}