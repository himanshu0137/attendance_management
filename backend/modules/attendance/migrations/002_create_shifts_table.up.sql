-- +goose Up
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

CREATE INDEX IF NOT EXISTS idx_shifts_tenant_id ON attendance.shifts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_shifts_active ON attendance.shifts(is_active);

-- +goose Down
DROP INDEX IF EXISTS idx_shifts_tenant_id;
DROP INDEX IF EXISTS idx_shifts_active;

DROP TABLE IF EXISTS attendance.shifts CASCADE;