-- +goose Up
CREATE TABLE IF NOT EXISTS org.employee_documents (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL,
    document_name TEXT NOT NULL,
    document_type TEXT NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (employee_id) REFERENCES org.employees(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_employee_documents_employee_id ON org.employee_documents(employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_documents_document_type ON org.employee_documents(document_type);

-- +goose Down
DROP INDEX IF EXISTS idx_employee_documents_employee_id;
DROP INDEX IF EXISTS idx_employee_documents_document_type;

DROP TABLE IF EXISTS org.employee_documents CASCADE;