-- Create user schema
CREATE SCHEMA IF NOT EXISTS user;

-- Create users table
CREATE TABLE IF NOT EXISTS user.users (
    tenant_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    keycloak_sub VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    role VARCHAR(50) NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- Ensure tenant_id and keycloak_sub combination is unique
    UNIQUE(tenant_id, keycloak_sub)
);

-- Create preferences table
CREATE TABLE IF NOT EXISTS user.preferences (
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