BEGIN;

CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    code TEXT NOT NULL, -- dtv, work_permit, ninety_day_report
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, code)
);


CREATE INDEX IF NOT EXISTS idx_services_tenant_id ON services(tenant_id);
CREATE INDEX IF NOT EXISTS idx_services_is_active ON services(is_active);


CREATE TABLE IF NOT EXISTS service_versions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    version_no INTEGER NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('draft', 'published', 'archived')),
    effective_from TIMESTAMPTZ,
    effective_to TIMESTAMPTZ,
    notes TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (service_id, version_no)
);


CREATE INDEX IF NOT EXISTS idx_service_versions_tenant_id ON service_versions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_service_versions_service_id ON service_versions(service_id);
CREATE INDEX IF NOT EXISTS idx_service_versions_status ON service_versions(status);


CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    code TEXT NOT NULL, -- passport, bank_statement, photo
    name TEXT NOT NULL,
    description TEXT,
    allowed_mime_types JSONB, -- ["application/pdf", "image/jpeg"]
    min_files INTEGER NOT NULL DEFAULT 1 CHECK (min_files >= 0),
    max_files INTEGER NOT NULL DEFAULT 1 CHECK (max_files >= min_files),
    max_file_size_bytes BIGINT CHECK (max_file_size_bytes IS NULL OR max_file_size_bytes > 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, code)
);


CREATE INDEX IF NOT EXISTS idx_documents_tenant_id ON documents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_documents_is_active ON documents(is_active);


CREATE TABLE IF NOT EXISTS service_documents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    service_version_id UUID NOT NULL REFERENCES service_versions(id) ON DELETE CASCADE,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE RESTRICT,
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    min_files_override INTEGER CHECK (min_files_override IS NULL OR min_files_override >= 0),
    max_files_override INTEGER CHECK (max_files_override IS NULL OR max_files_override >= 0),
    display_order INTEGER NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (service_version_id, document_id),
    CHECK (
        min_files_override IS NULL
        OR max_files_override IS NULL
        OR min_files_override <= max_files_override
    )
);


CREATE INDEX IF NOT EXISTS idx_service_documents_tenant_id ON service_documents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_service_documents_service_version_id
ON service_documents(service_version_id);
CREATE INDEX IF NOT EXISTS idx_service_documents_document_id
ON service_documents(document_id);


CREATE TABLE IF NOT EXISTS submission_countries (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    country_code TEXT NOT NULL,
    country_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, country_code)
);


CREATE INDEX IF NOT EXISTS idx_submission_countries_tenant_id
ON submission_countries(tenant_id);
CREATE INDEX IF NOT EXISTS idx_submission_countries_is_active
ON submission_countries(is_active);


CREATE TABLE IF NOT EXISTS service_applications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    service_version_id UUID NOT NULL REFERENCES service_versions(id) ON DELETE RESTRICT,
    status TEXT NOT NULL CHECK (
        status IN ('draft', 'in_progress', 'submitted', 'waiting_for_client', 'under_review', 'approved', 'rejected', 'cancelled')
    ),
    submitted_at TIMESTAMPTZ,
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_service_applications_tenant_id ON service_applications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_service_applications_user_id ON service_applications(user_id);
CREATE INDEX IF NOT EXISTS idx_service_applications_service_version_id ON service_applications(service_version_id);
CREATE INDEX IF NOT EXISTS idx_service_applications_status ON service_applications(status);
CREATE INDEX IF NOT EXISTS idx_service_applications_assignee_id
ON service_applications(assignee_id);
CREATE INDEX IF NOT EXISTS idx_service_applications_tenant_id_status_created_at
ON service_applications(tenant_id, status, created_at);
CREATE INDEX IF NOT EXISTS idx_service_applications_tenant_id_assignee_id_status_created_at
ON service_applications(tenant_id, assignee_id, status, created_at);
CREATE INDEX IF NOT EXISTS idx_service_applications_tenant_id_service_version_id_status_created_at
ON service_applications(tenant_id, service_version_id, status, created_at);


CREATE TABLE IF NOT EXISTS dtv_service_applications (
    service_application_id UUID PRIMARY KEY REFERENCES service_applications(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    submission_country_id UUID REFERENCES submission_countries(id) ON DELETE RESTRICT,
    submission_country_arrival_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_dtv_service_applications_tenant_id
ON dtv_service_applications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_dtv_service_applications_submission_country_id
ON dtv_service_applications(submission_country_id);
CREATE INDEX IF NOT EXISTS idx_dtv_service_applications_submission_country_arrival_date
ON dtv_service_applications(submission_country_arrival_date);
CREATE INDEX IF NOT EXISTS idx_dtv_service_applications_tenant_id_submission_country_id_submission_country_arrival_date
ON dtv_service_applications(tenant_id, submission_country_id, submission_country_arrival_date);


CREATE TABLE IF NOT EXISTS document_submissions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    service_application_id UUID NOT NULL REFERENCES service_applications(id) ON DELETE CASCADE,
    service_document_id UUID NOT NULL REFERENCES service_documents(id) ON DELETE RESTRICT,
    status TEXT NOT NULL CHECK (status IN ('pending', 'needs_changes', 'approved', 'rejected')),
    submitted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    submitted_at TIMESTAMPTZ,
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    feedback TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (service_application_id, service_document_id)
);


CREATE INDEX IF NOT EXISTS idx_document_submissions_tenant_id
ON document_submissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_document_submissions_service_application_id
ON document_submissions(service_application_id);
CREATE INDEX IF NOT EXISTS idx_document_submissions_service_document_id
ON document_submissions(service_document_id);
CREATE INDEX IF NOT EXISTS idx_document_submissions_status
ON document_submissions(status);


CREATE TABLE IF NOT EXISTS document_files (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    document_submission_id UUID NOT NULL REFERENCES document_submissions(id) ON DELETE CASCADE,
    storage_path TEXT NOT NULL,
    file_name TEXT,
    mime_type TEXT,
    file_size_bytes BIGINT CHECK (file_size_bytes IS NULL OR file_size_bytes >= 0),
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);


CREATE INDEX IF NOT EXISTS idx_document_files_tenant_id
ON document_files(tenant_id);
CREATE INDEX IF NOT EXISTS idx_document_files_document_submission_id
ON document_files(document_submission_id);
CREATE INDEX IF NOT EXISTS idx_document_files_is_active
ON document_files(is_active);


-- V1 step model:
-- - step code list is managed in app logic (for example: begin, upload, staff_review, request_changes)
-- - rows are append-only style events for service application workflow progression
CREATE TABLE IF NOT EXISTS application_steps (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    service_application_id UUID NOT NULL REFERENCES service_applications(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_application_steps_tenant_id
ON application_steps(tenant_id);
CREATE INDEX IF NOT EXISTS idx_application_steps_service_application_id
ON application_steps(service_application_id);
CREATE INDEX IF NOT EXISTS idx_application_steps_code
ON application_steps(code);
CREATE INDEX IF NOT EXISTS idx_application_steps_created_at
ON application_steps(created_at);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_application_steps_service_application_id_code
ON application_steps(service_application_id, code);
CREATE INDEX IF NOT EXISTS idx_application_steps_service_application_id_created_at
ON application_steps(service_application_id, created_at);
CREATE INDEX IF NOT EXISTS idx_application_steps_tenant_id_service_application_id_created_at
ON application_steps(tenant_id, service_application_id, created_at);


END;