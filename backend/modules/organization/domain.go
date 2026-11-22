package organization

import "time"

// AllowedPunchMethods is the set of punch methods a tenant may enable.
const (
	// FACESCAN_KIOSK is the kiosk face recognition biometric.
	FACESCAN_KIOSK = "FACESCAN_KIOSK"
	// FINGERSCAN_KIOSK is the kiosk finger scan biometric.
	FINGERSCAN_KIOSK = "FINGERSCAN_KIOSK"
	// WEB_LOGIN is the Angular web portal login.
	WEB_LOGIN = "WEB_LOGIN"
	// MOBILE_LOGIN is the Flutter employee app login.
	MOBILE_LOGIN = "MOBILE_LOGIN"
)

// OrgPolicy is the per-tenant attendance policy.
type OrgPolicy struct {
	// Methods is the list of enabled punch methods.
	Methods []string `json:"methods"`
	// SalarySlipEnabled controls whether payroll may generate salary slips.
	SalarySlipEnabled bool `json:"salary_slip_enabled"`
	// ShiftStart/ShiftEnd bound the working day in local time.
	ShiftStart time.Time `json:"shift_start_time"`
	ShiftEnd   time.Time `json:"shift_end_time"`
}

// IsEnabled reports whether the given punch method is enabled for the policy.
func (p OrgPolicy) IsEnabled(method string) bool {
	for _, m := range p.Methods {
		if m == method {
			return true
		}
	}
	return false
}

// Employee is a Keycloak-backed org user.
type Employee struct {
	ID          string `json:"id"`
	KeycloakSub string `json:"keycloak_sub"`
	Role        string `json:"role"`
}
