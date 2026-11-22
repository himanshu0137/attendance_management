package payroll

import (
	"math/big"
	"time"
)

// Money is an exact integer-based currency amount in minor units.
type Money struct {
	amount *big.Int
}

// NewMoney parses a decimal string into Money.
func NewMoney(s string) Money {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("payroll: invalid money")
	}
	return Money{amount: v}
}

// Add returns the sum.
func (m Money) Add(n Money) Money {
	return Money{amount: new(big.Int).Add(m.amount, n.amount)}
}

// Sub returns the difference.
func (m Money) Sub(n Money) Money {
	return Money{amount: new(big.Int).Sub(m.amount, n.amount)}
}

// Div returns the quotient, truncating towards zero.
func (m Money) Div(n Money) Money {
	if n.amount.Sign() == 0 {
		return Money{}
	}
	return Money{amount: new(big.Int).Div(m.amount, n.amount)}
}

// String renders the amount as a decimal string.
func (m Money) String() string {
	return m.amount.String()
}

// Equal reports whether both amounts are the same.
func (m Money) Equal(n Money) bool {
	return m.amount.Cmp(n.amount) == 0
}

// Cmp compares this amount to another.
func (m Money) Cmp(n Money) int {
	return m.amount.Cmp(n.amount)
}

// TenantSettings is the tenant's salary slip eligibility.
type TenantSettings struct {
	SalarySlipEnabled bool `json:"salary_slip_enabled"`
}

// SalaryStructure is the per-employee monthly salary and payable day count.
type SalaryStructure struct {
	EmployeeID   string   `json:"employee_id"`
	BaseSalary   Money    `json:"base_salary"`
	DaysPerMonth int64    `json:"days_per_month"`
}

// Day is one payable day in a slip.
type Day struct {
	Date    time.Time `json:"date"`
	Minutes int64     `json:"minutes"`
}

// SalarySlip is a generated monthly salary slip.
type SalarySlip struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	EmployeeID string    `json:"employee_id"`
	Period     string    `json:"period"`
	BaseSalary Money     `json:"base_salary"`
	Days       []Day     `json:"days"`
}
