-- +goose Up
CREATE SCHEMA IF NOT EXISTS payroll;

CREATE TABLE IF NOT EXISTS payroll.tenant_settings (
	tenant_id TEXT NOT NULL PRIMARY KEY,
	salary_slip_enabled BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS payroll.salary_structures (
	tenant_id TEXT NOT NULL,
	employee_id TEXT NOT NULL,
	base_salary NUMERIC(18,2) NOT NULL,
	days_per_month INT NOT NULL DEFAULT 0,
	PRIMARY KEY (tenant_id, employee_id),
	FOREIGN KEY (tenant_id) REFERENCES payroll.tenant_settings(tenant_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payroll.salary_slips (
	id TEXT NOT NULL PRIMARY KEY,
	tenant_id TEXT NOT NULL,
	employee_id TEXT NOT NULL,
	period TEXT NOT NULL,
	base_salary NUMERIC(18,2) NOT NULL,
	days JSONB NOT NULL DEFAULT '[]'::jsonb,
	FOREIGN KEY (tenant_id) REFERENCES payroll.tenant_settings(tenant_id) ON DELETE CASCADE,
	FOREIGN KEY (employee_id) REFERENCES payroll.salary_structures(employee_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS payroll.salary_slips;
DROP TABLE IF EXISTS payroll.salary_structures;
DROP TABLE IF EXISTS payroll.tenant_settings;
