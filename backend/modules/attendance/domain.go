package attendance

import (
	"errors"
	"time"
)

// PunchMethod is the biometric or login channel used for a punch.
type PunchMethod string

const (
	FacescanKiosk   PunchMethod = "FACESCAN_KIOSK"
	FingerScanKiosk PunchMethod = "FINGERSCAN_KIOSK"
	WebLogin        PunchMethod = "WEB_LOGIN"
	MobileLogin     PunchMethod = "MOBILE_LOGIN"
)

// TenantPolicy is the per-tenant attendance policy held locally.
type TenantPolicy struct {
	Methods    []string  `json:"methods"`
	ShiftStart time.Time `json:"shift_start_time"`
	ShiftEnd   time.Time `json:"shift_end_time"`
}

// IsEnabled reports whether the method is enabled for the policy.
func (p TenantPolicy) IsEnabled(method PunchMethod) bool {
	for _, m := range p.Methods {
		if m == string(method) {
			return true
		}
	}
	return false
}

// ErrPolicyNotFound is returned when a tenant has no local policy yet.
type ErrPolicyNotFound struct{ tenantID string }

func (e ErrPolicyNotFound) Error() string {
	return "attendance: policy not found for tenant " + e.tenantID
}

// ErrMethodDisabled is returned when the punch method is not enabled.
var ErrMethodDisabled = errors.New("attendance: punch method not enabled for this tenant")

// ErrOutsideShift is returned when the punch is outside the configured shift.
var ErrOutsideShift = errors.New("attendance: punch outside the configured shift")

// PunchLog is a raw punch event.
type PunchLog struct {
	ID         int64       `json:"id"`
	TenantID   string      `json:"tenant_id"`
	EmployeeID string      `json:"employee_id"`
	Method     PunchMethod `json:"method"`
	Start      time.Time   `json:"start"`
	End        *time.Time  `json:"end"`
}

// Report is the aggregated attendance for a date range.
type Report struct {
	PresentDays   int64 `json:"present_days"`
	WorkedMinutes int64 `json:"worked_minutes"`
}

// Shift represents an employee shift
type Shift struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AttendanceRule represents an attendance policy rule
type AttendanceRule struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Name        string                 `json:"rule_name"`
	Description string                 `json:"description"`
	RuleType    string                 `json:"rule_type"` // e.g., "late_arrival", "early_departure", "overtime"
	Condition   map[string]interface{} `json:"condition"` // flexible conditions for the rule
	Action      map[string]interface{} `json:"action"`    // action to take when rule is triggered
	IsActive    bool                   `json:"is_active"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
