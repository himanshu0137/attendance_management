-- Single schema for all tables
CREATE SCHEMA IF NOT EXISTS attendance;

-- Create ENUM type for punch methods
CREATE TYPE IF NOT EXISTS attendance.punch_method AS ENUM ('FACESCAN_KIOSK', 'FINGERSCAN_KIOSK', 'WEB_LOGIN', 'MOBILE_LOGIN');

-- Tenant policies table
CREATE TABLE IF NOT EXISTS attendance.tenant_policies (
	tenant_id TEXT NOT NULL PRIMARY KEY,
	methods attendance.punch_method[] NOT NULL DEFAULT '{}'::attendance.punch_method[],
	shift_start_time TIME NOT NULL DEFAULT '09:00'::time,
	shift_end_time TIME NOT NULL DEFAULT '18:00'::time
);

-- Attendance logs table
CREATE TABLE IF NOT EXISTS attendance.logs (
	id BIGSERIAL PRIMARY KEY,
	tenant_id TEXT NOT NULL,
	employee_id TEXT NOT NULL,
	method attendance.punch_method NOT NULL,
	start_time TIME NOT NULL,
	end_time TIME,
	FOREIGN KEY (tenant_id) REFERENCES attendance.tenant_policies(tenant_id) ON DELETE CASCADE
);

-- Shifts table
CREATE TABLE IF NOT EXISTS attendance.shifts (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (tenant_id) REFERENCES attendance.tenant_policies(tenant_id) ON DELETE CASCADE
);

-- Create indexes for shifts
CREATE INDEX IF NOT EXISTS idx_shifts_tenant_id ON attendance.shifts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_shifts_active ON attendance.shifts(is_active);

-- Leave requests table
CREATE TABLE IF NOT EXISTS attendance.leave_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    leave_type VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    approved_by VARCHAR(255),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for leave requests
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON attendance.leave_requests(employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_tenant_id ON attendance.leave_requests(tenant_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON attendance.leave_requests(status);
CREATE INDEX IF NOT EXISTS idx_leave_requests_start_date ON attendance.leave_requests(start_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_end_date ON attendance.leave_requests(end_date);

-- Add check constraint to ensure end_date is not before start_date
ALTER TABLE attendance.leave_requests 
ADD CONSTRAINT chk_end_date_after_start_date 
CHECK (end_date >= start_date);

-- Attendance rules table
CREATE TABLE IF NOT EXISTS attendance.attendance_rules (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    rule_name TEXT NOT NULL,
    description TEXT,
    rule_type TEXT NOT NULL, -- e.g., "late_arrival", "early_departure", "overtime"
    condition JSONB NOT NULL, -- flexible conditions for the rule
    action JSONB NOT NULL, -- action to take when rule is triggered
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (tenant_id) REFERENCES attendance.tenant_policies(tenant_id) ON DELETE CASCADE
);

-- Create indexes for attendance rules
CREATE INDEX IF NOT EXISTS idx_attendance_rules_tenant_id ON attendance.attendance_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_attendance_rules_type ON attendance.attendance_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_attendance_rules_active ON attendance.attendance_rules(is_active);

-- Organization tables - consolidated into single schema
CREATE TABLE IF NOT EXISTS attendance.organization_policies (
	tenant_id TEXT NOT NULL PRIMARY KEY,
	methods TEXT[] NOT NULL DEFAULT '{}'::text[],
	salary_slip_enabled BOOLEAN NOT NULL DEFAULT false,
	shift_start_time TIME NOT NULL DEFAULT '09:00'::time,
	shift_end_time TIME NOT NULL DEFAULT '18:00'::time
);

CREATE TABLE IF NOT EXISTS attendance.employees (
	tenant_id TEXT NOT NULL,
	keycloak_id UUID NOT NULL,
	role TEXT NOT NULL,
	first_name TEXT,
	last_name TEXT,
	email TEXT,
	employee_code TEXT UNIQUE,
	designation TEXT,
	department TEXT,
	manager_id UUID,
	branch_id UUID,
	joining_date DATE,
	status TEXT DEFAULT 'active',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	PRIMARY KEY (tenant_id, keycloak_id),
	FOREIGN KEY (tenant_id) REFERENCES attendance.organization_policies(tenant_id) ON DELETE CASCADE
);

-- Add foreign key constraint for manager
ALTER TABLE attendance.employees 
ADD CONSTRAINT fk_employees_manager 
FOREIGN KEY (manager_id) REFERENCES attendance.employees(id) ON DELETE SET NULL;

-- Create indexes for employees
CREATE INDEX IF NOT EXISTS idx_employees_employee_code ON attendance.employees(employee_code);
CREATE INDEX IF NOT EXISTS idx_employees_designation ON attendance.employees(designation);
CREATE INDEX IF NOT EXISTS idx_employees_department ON attendance.employees(department);
CREATE INDEX IF NOT EXISTS idx_employees_status ON attendance.employees(status);
CREATE INDEX IF NOT EXISTS idx_employees_joining_date ON attendance.employees(joining_date);

-- Salary structures table
CREATE TABLE IF NOT EXISTS attendance.salary_structures (
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
    
    FOREIGN KEY (employee_id) REFERENCES attendance.employees(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES attendance.organization_policies(tenant_id) ON DELETE CASCADE
);

-- Create indexes for salary structures
CREATE INDEX IF NOT EXISTS idx_salary_structures_employee_id ON attendance.salary_structures(employee_id);
CREATE INDEX IF NOT EXISTS idx_salary_structures_tenant_id ON attendance.salary_structures(tenant_id);
CREATE INDEX IF NOT EXISTS idx_salary_structures_effective_date ON attendance.salary_structures(effective_date);

-- Employee documents table
CREATE TABLE IF NOT EXISTS attendance.employee_documents (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    document_name TEXT NOT NULL,
    document_type TEXT NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (employee_id) REFERENCES attendance.employees(id) ON DELETE CASCADE
);

-- Create indexes for employee documents
CREATE INDEX IF NOT EXISTS idx_employee_documents_employee_id ON attendance.employee_documents(employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_documents_document_type ON attendance.employee_documents(document_type);

-- Payroll tables - consolidated into single schema
CREATE TABLE IF NOT EXISTS attendance.payroll_settings (
	tenant_id TEXT NOT NULL PRIMARY KEY,
	salary_slip_enabled BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS attendance.payroll_salary_structures (
	tenant_id TEXT NOT NULL,
	employee_id TEXT NOT NULL,
	base_salary NUMERIC(18,2) NOT NULL,
	days_per_month INT NOT NULL DEFAULT 0,
	PRIMARY KEY (tenant_id, employee_id),
	FOREIGN KEY (tenant_id) REFERENCES attendance.payroll_settings(tenant_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS attendance.payroll_salary_slips (
	id TEXT NOT NULL PRIMARY KEY,
	tenant_id TEXT NOT NULL,
	employee_id TEXT NOT NULL,
	period TEXT NOT NULL,
	base_salary NUMERIC(18,2) NOT NULL,
	days JSONB NOT NULL DEFAULT '[]'::jsonb,
	FOREIGN KEY (tenant_id) REFERENCES attendance.payroll_settings(tenant_id) ON DELETE CASCADE,
	FOREIGN KEY (employee_id) REFERENCES attendance.payroll_salary_structures(employee_id) ON DELETE CASCADE
);

-- Statutory compliance table
CREATE TABLE IF NOT EXISTS attendance.statutory_compliance (
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
    
    FOREIGN KEY (tenant_id) REFERENCES attendance.payroll_settings(tenant_id) ON DELETE CASCADE
);

-- Create indexes for statutory compliance
CREATE INDEX IF NOT EXISTS idx_statutory_compliance_tenant_id ON attendance.statutory_compliance(tenant_id);
CREATE INDEX IF NOT EXISTS idx_statutory_compliance_type ON attendance.statutory_compliance(compliance_type);
CREATE INDEX IF NOT EXISTS idx_statutory_compliance_active ON attendance.statutory_compliance(is_active);

-- Employee statutory deductions table
CREATE TABLE IF NOT EXISTS attendance.employee_statutory_deductions (
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
    
    FOREIGN KEY (employee_id) REFERENCES attendance.payroll_salary_structures(employee_id) ON DELETE CASCADE,
    FOREIGN KEY (statutory_id) REFERENCES attendance.statutory_compliance(id) ON DELETE CASCADE
);

-- Create indexes for employee statutory deductions
CREATE INDEX IF NOT EXISTS idx_employee_statutory_deductions_employee_id ON attendance.employee_statutory_deductions(employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_statutory_deductions_statutory_id ON attendance.employee_statutory_deductions(statutory_id);

-- User tables - consolidated into single schema
CREATE TABLE IF NOT EXISTS attendance.users (
    tenant_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keycloak_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    role VARCHAR(50) NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Ensure tenant_id and keycloak_id combination is unique
    UNIQUE(tenant_id, keycloak_id)
);

-- Preferences table
CREATE TABLE IF NOT EXISTS attendance.user_preferences (
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    theme VARCHAR(20) DEFAULT 'light',
    notifications BOOLEAN DEFAULT TRUE,
    language VARCHAR(10) DEFAULT 'en',
    date_format VARCHAR(20) DEFAULT 'YYYY-MM-DD',
    auto_clock_in BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (tenant_id, user_id)
);