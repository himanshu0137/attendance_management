package leave

import (
	"context"
	"testing"
	"time"

	"uuid"

	"github.com/stretchr/testify/assert"
)

// Mock repository for testing
type mockRepository struct {
	leaveRequests map[uuid.UUID]*LeaveRequest
}

func (m *mockRepository) CreateLeaveRequest(ctx context.Context, leave LeaveRequest) error {
	if m.leaveRequests == nil {
		m.leaveRequests = make(map[uuid.UUID]*LeaveRequest)
	}
	m.leaveRequests[leave.ID] = &leave
	return nil
}

func (m *mockRepository) GetLeaveRequest(ctx context.Context, id uuid.UUID) (*LeaveRequest, error) {
	if leave, exists := m.leaveRequests[id]; exists {
		return leave, nil
	}
	return nil, nil
}

func (m *mockRepository) ListLeaveRequests(ctx context.Context, tenantID, status string, page, limit int) ([]*LeaveRequest, int, error) {
	var leaves []*LeaveRequest
	for _, leave := range m.leaveRequests {
		if leave.TenantID == tenantID && (status == "" || leave.Status == status) {
			leaves = append(leaves, leave)
		}
	}
	return leaves, len(leaves), nil
}

func (m *mockRepository) UpdateLeaveStatus(ctx context.Context, id uuid.UUID, status string, approvedBy *string, approvedAt *time.Time) error {
	if leave, exists := m.leaveRequests[id]; exists {
		leave.Status = status
		leave.ApprovedBy = approvedBy
		leave.ApprovedAt = approvedAt
		leave.UpdatedAt = time.Now()
		return nil
	}
	return nil
}

func (m *mockRepository) ListMyLeaveRequests(ctx context.Context, employeeID string, page, limit int) ([]*LeaveRequest, int, error) {
	var leaves []*LeaveRequest
	for _, leave := range m.leaveRequests {
		if leave.EmployeeID == employeeID {
			leaves = append(leaves, leave)
		}
	}
	return leaves, len(leaves), nil
}

func TestCreateLeaveRequest(t *testing.T) {
	repo := &mockRepository{}
	service := NewService(repo)

	ctx := context.Background()
	req := CreateLeaveRequest{
		EmployeeID:  "emp-123",
		TenantID:    "tenant-456",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		LeaveType:   "Casual",
		Description: "Test leave request",
	}

	leave, err := service.CreateLeaveRequest(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, leave)
	assert.Equal(t, "PENDING", leave.Status)
	assert.Equal(t, req.EmployeeID, leave.EmployeeID)
	assert.Equal(t, req.TenantID, leave.TenantID)
	assert.Equal(t, req.LeaveType, leave.LeaveType)
}

func TestApproveLeaveRequest(t *testing.T) {
	repo := &mockRepository{}
	service := NewService(repo)

	ctx := context.Background()

	// First create a leave request
	req := CreateLeaveRequest{
		EmployeeID:  "emp-123",
		TenantID:    "tenant-456",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		LeaveType:   "Casual",
		Description: "Test leave request",
	}

	leave, err := service.CreateLeaveRequest(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, leave)

	// Then approve it
	approvedBy := "manager-789"
	err = service.ApproveLeaveRequest(ctx, leave.ID, approvedBy)
	assert.NoError(t, err)

	// Verify approval
	updatedLeave, err := service.GetLeaveRequest(ctx, leave.ID)
	assert.NoError(t, err)
	assert.Equal(t, "APPROVED", updatedLeave.Status)
	assert.Equal(t, &approvedBy, updatedLeave.ApprovedBy)
}

func TestRejectLeaveRequest(t *testing.T) {
	repo := &mockRepository{}
	service := NewService(repo)

	ctx := context.Background()

	// First create a leave request
	req := CreateLeaveRequest{
		EmployeeID:  "emp-123",
		TenantID:    "tenant-456",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		LeaveType:   "Casual",
		Description: "Test leave request",
	}

	leave, err := service.CreateLeaveRequest(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, leave)

	// Then reject it
	rejectedBy := "manager-789"
	err = service.RejectLeaveRequest(ctx, leave.ID, rejectedBy)
	assert.NoError(t, err)

	// Verify rejection
	updatedLeave, err := service.GetLeaveRequest(ctx, leave.ID)
	assert.NoError(t, err)
	assert.Equal(t, "REJECTED", updatedLeave.Status)
	assert.Equal(t, &rejectedBy, updatedLeave.ApprovedBy)
}

func TestListMyLeaveRequests(t *testing.T) {
	repo := &mockRepository{}
	service := NewService(repo)

	ctx := context.Background()

	// Create multiple leave requests for same employee
	req1 := CreateLeaveRequest{
		EmployeeID:  "emp-123",
		TenantID:    "tenant-456",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		LeaveType:   "Casual",
		Description: "Test leave request 1",
	}

	req2 := CreateLeaveRequest{
		EmployeeID:  "emp-123",
		TenantID:    "tenant-456",
		StartDate:   time.Now().Add(48 * time.Hour),
		EndDate:     time.Now().Add(72 * time.Hour),
		LeaveType:   "Sick",
		Description: "Test leave request 2",
	}

	// Create requests
	leave1, err := service.CreateLeaveRequest(ctx, req1)
	assert.NoError(t, err)

	leave2, err := service.CreateLeaveRequest(ctx, req2)
	assert.NoError(t, err)

	// List employee's leave requests
	leaves, total, err := service.ListMyLeaveRequests(ctx, "emp-123", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, leaves, 2)

	// Verify both are returned
	foundIDs := make(map[uuid.UUID]bool)
	for _, leave := range leaves {
		foundIDs[leave.ID] = true
	}
	assert.True(t, foundIDs[leave1.ID])
	assert.True(t, foundIDs[leave2.ID])
}
