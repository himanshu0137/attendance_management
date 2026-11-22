-- +goose Up
CREATE SCHEMA IF NOT EXISTS attendance;

CREATE TYPE IF NOT EXISTS punch_method AS ENUM ('FACESCAN_KIOSK', 'FINGERSCAN_KIOSK', 'WEB_LOGIN', 'MOBILE_LOGIN');

CREATE TABLE IF NOT EXISTS attendance.tenant_policies (
	tenant_id TEXT NOT NULL PRIMARY KEY,
	methods punch_method[] NOT NULL DEFAULT '{}'::punch_method[],
	shift_start_time TIME NOT NULL DEFAULT '09:00'::time,
	shift_end_time TIME NOT NULL DEFAULT '18:00'::time
);

CREATE TABLE IF NOT EXISTS attendance.logs (
	id BIGSERIAL PRIMARY KEY,
	tenant_id TEXT NOT NULL,
	employee_id TEXT NOT NULL,
	method punch_method NOT NULL,
	start_time TIME NOT NULL,
	end_time TIME,
	FOREIGN KEY (tenant_id) REFERENCES attendance.tenant_policies(tenant_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS attendance.logs;
DROP TABLE IF EXISTS attendance.tenant_policies;
DROP TYPE IF EXISTS punch_method;
