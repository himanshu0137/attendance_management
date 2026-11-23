-- +goose Up
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

CREATE INDEX IF NOT EXISTS idx_attendance_rules_tenant_id ON attendance.attendance_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_attendance_rules_type ON attendance.attendance_rules(rule_type);
CREATE INDEX IF NOT EXISTS idx_attendance_rules_active ON attendance.attendance_rules(is_active);

-- +goose Down
DROP INDEX IF EXISTS idx_attendance_rules_tenant_id;
DROP INDEX IF EXISTS idx_attendance_rules_type;
DROP INDEX IF EXISTS idx_attendance_rules_active;

DROP TABLE IF EXISTS attendance.attendance_rules CASCADE;