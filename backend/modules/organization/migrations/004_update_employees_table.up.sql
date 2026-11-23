-- +goose Up
ALTER TABLE org.employees 
ADD COLUMN IF NOT EXISTS first_name TEXT,
ADD COLUMN IF NOT EXISTS last_name TEXT,
ADD COLUMN IF NOT EXISTS email TEXT,
ADD COLUMN IF NOT EXISTS employee_code TEXT UNIQUE,
ADD COLUMN IF NOT EXISTS designation TEXT,
ADD COLUMN IF NOT EXISTS department TEXT,
ADD COLUMN IF NOT EXISTS manager_id UUID,
ADD COLUMN IF NOT EXISTS branch_id UUID,
ADD COLUMN IF NOT EXISTS joining_date DATE,
ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active',
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT NOW(),
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();

-- Add foreign key constraints
ALTER TABLE org.employees 
ADD CONSTRAINT fk_employees_manager 
FOREIGN KEY (manager_id) REFERENCES org.employees(id) ON DELETE SET NULL;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_employees_employee_code ON org.employees(employee_code);
CREATE INDEX IF NOT EXISTS idx_employees_designation ON org.employees(designation);
CREATE INDEX IF NOT EXISTS idx_employees_department ON org.employees(department);
CREATE INDEX IF NOT EXISTS idx_employees_status ON org.employees(status);
CREATE INDEX IF NOT EXISTS idx_employees_joining_date ON org.employees(joining_date);

-- +goose Down
DROP INDEX IF EXISTS idx_employees_employee_code;
DROP INDEX IF EXISTS idx_employees_designation;
DROP INDEX IF EXISTS idx_employees_department;
DROP INDEX IF EXISTS idx_employees_status;
DROP INDEX IF EXISTS idx_employees_joining_date;

ALTER TABLE org.employees 
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name,
DROP COLUMN IF EXISTS email,
DROP COLUMN IF EXISTS employee_code,
DROP COLUMN IF EXISTS designation,
DROP COLUMN IF EXISTS department,
DROP COLUMN IF EXISTS manager_id,
DROP COLUMN IF EXISTS branch_id,
DROP COLUMN IF EXISTS joining_date,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE org.employees 
DROP CONSTRAINT IF EXISTS fk_employees_manager;