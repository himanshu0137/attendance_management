package attendance

import (
	"errors"
	"time"
)

// PunchMethod is the biometric or login channel used for a punch.
type PunchMethod string

const (
	FacescanKiosk PunchMethod = "FACESCAN_KIOSK"
	FingerScanKiosk PunchMethod = "FINGERSCAN_KIOSK"
	WebLogin      PunchMethod = "WEB_LOGIN"
	MobileLogin   PunchMethod = "MOBILE_LOGIN"
)

// TenantPolicy is the per-tenant attendance policy held locally.
type TenantPolicy struct {
	Methods   []string `json:"methods"`
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
	ID       int64     `json:"id"`
	TenantID string    `json:"tenant_id"`
	EmployeeID string   `json:"employee_id"`
	Method   PunchMethod `json:"method"`
	Start    time.Time `json:"start"`
	End      *time.Time `json:"end"`
}

// Report is the aggregated attendance for a date range.
type Report struct {
	PresentDays   int64 `json:"present_days"`
	WorkedMinutes int64 `json:"worked_minutes"`
}
