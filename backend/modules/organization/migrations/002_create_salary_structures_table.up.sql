-- +goose Up
CREATE TABLE IF NOT EXISTS org.salary_structures (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    tenant_id TEXT NOT NULL,
    pay_basis TEXT NOT NULL, -- "monthly" or "daily"
    basic NUMERIC(15,2) NOT NULL DEFAULT 0,
    hra NUMERIC(15,2) NOT NULL DEFAULT 0,
    da NUMERIC(15,2) NOT NULL DEFAULT 0,
    conveyance NUMERIC(15,2) NOT NULL DEFAULT 0,
    special_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    medical_allowance NUMERIC(15,2) NOT NULL DEFAULT 0,
    bonus NUMERIC(15,2) NOT NULL DEFAULT 0,
    overtime_rate NUMERIC(15,2) NOT NULL DEFAULT 0,
    effective_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (employee_id) REFERENCES org.employees(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES org.policies(tenant_id) ON DELETE CASCADE
);

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_salary_structures_employee_id ON org.salary_structures(employee_id);
CREATE INDEX IF NOT EXISTS idx_salary_structures_tenant_id ON org.salary_structures(tenant_id);
CREATE INDEX IF NOT EXISTS idx_salary_structures_effective_date ON org.salary_structures(effective_date);

-- +goose Down
DROP INDEX IF EXISTS idx_salary_structures_employee_id;
DROP INDEX IF EXISTS idx_salary_structures_tenant_id;
DROP INDEX IF EXISTS idx_salary_structures_effective_date;

DROP TABLE IF EXISTS org.salary_structures CASCADE;