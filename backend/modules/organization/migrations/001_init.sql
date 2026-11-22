-- +goose Up
CREATE SCHEMA IF NOT EXISTS org;

CREATE TABLE IF NOT EXISTS org.policies (
	tenant_id TEXT NOT NULL PRIMARY KEY,
	methods TEXT[] NOT NULL DEFAULT '{}'::text[],
	salary_slip_enabled BOOLEAN NOT NULL DEFAULT false,
	shift_start_time TIME NOT NULL DEFAULT '09:00'::time,
	shift_end_time TIME NOT NULL DEFAULT '18:00'::time
);

CREATE TABLE IF NOT EXISTS org.employees (
	tenant_id TEXT NOT NULL,
	keycloak_sub UUID NOT NULL,
	role TEXT NOT NULL,
	PRIMARY KEY (tenant_id, keycloak_sub),
	FOREIGN KEY (tenant_id) REFERENCES org.policies(tenant_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS org.employees CASCADE;
DROP TABLE IF EXISTS org.policies;
