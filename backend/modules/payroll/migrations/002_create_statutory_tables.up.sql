-- +goose Up
CREATE TABLE IF NOT EXISTS payroll.statutory_compliance (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    compliance_type TEXT NOT NULL, -- e.g., "ESIC", "EPF", "Professional_Tax", "Income_Tax"
    applicable_from DATE NOT NULL,
    applicable_to DATE, -- NULL means still applicable
    rate NUMERIC(18,2) NOT NULL,
    max_limit NUMERIC(18,2), -- Optional maximum limit for calculation
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (tenant_id) REFERENCES payroll.tenant_settings(tenant_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_statutory_compliance_tenant_id ON payroll.statutory_compliance(tenant_id);
CREATE INDEX IF NOT EXISTS idx_statutory_compliance_type ON payroll.statutory_compliance(compliance_type);
CREATE INDEX IF NOT EXISTS idx_statutory_compliance_active ON payroll.statutory_compliance(is_active);

-- Create table for storing employee-specific statutory deductions
CREATE TABLE IF NOT EXISTS payroll.employee_statutory_deductions (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id TEXT NOT NULL,
    statutory_id UUID NOT NULL,
    applicable_from DATE NOT NULL,
    applicable_to DATE, -- NULL means still applicable
    rate NUMERIC(18,2) NOT NULL,
    max_limit NUMERIC(18,2), -- Optional maximum limit for calculation
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (employee_id) REFERENCES payroll.salary_structures(employee_id) ON DELETE CASCADE,
    FOREIGN KEY (statutory_id) REFERENCES payroll.statutory_compliance(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_employee_statutory_deductions_employee_id ON payroll.employee_statutory_deductions(employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_statutory_deductions_statutory_id ON payroll.employee_statutory_deductions(statutory_id);

-- +goose Down
DROP INDEX IF EXISTS idx_employee_statutory_deductions_employee_id;
DROP INDEX IF EXISTS idx_employee_statutory_deductions_statutory_id;
DROP TABLE IF EXISTS payroll.employee_statutory_deductions CASCADE;

DROP INDEX IF EXISTS idx_statutory_compliance_tenant_id;
DROP INDEX IF EXISTS idx_statutory_compliance_type;
DROP INDEX IF EXISTS idx_statutory_compliance_active;
DROP TABLE IF EXISTS payroll.statutory_compliance CASCADE;