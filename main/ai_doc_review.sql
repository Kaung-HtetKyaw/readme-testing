BEGIN;

CREATE TABLE IF NOT EXISTS ai_doc_reviews (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    service_application_id UUID NOT NULL REFERENCES service_applications(id) ON DELETE CASCADE,
    document_submission_id UUID NOT NULL REFERENCES document_submissions(id) ON DELETE CASCADE,
    verdict TEXT CHECK (verdict IN ('approved', 'rejected', 'unsure')),
    feedback TEXT,
    triggered_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    triggered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);


CREATE INDEX IF NOT EXISTS idx_ai_doc_reviews_tenant_id
ON ai_doc_reviews(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_doc_reviews_service_application_id
ON ai_doc_reviews(service_application_id);
CREATE INDEX IF NOT EXISTS idx_ai_doc_reviews_document_submission_id
ON ai_doc_reviews(document_submission_id);
CREATE INDEX IF NOT EXISTS idx_ai_doc_reviews_triggered_at
ON ai_doc_reviews(triggered_at);


CREATE TABLE IF NOT EXISTS ai_doc_review_run_files (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    ai_doc_review_id UUID NOT NULL REFERENCES ai_doc_reviews(id) ON DELETE CASCADE,
    document_file_id UUID NOT NULL REFERENCES document_files(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (ai_doc_review_id, document_file_id)
);


CREATE INDEX IF NOT EXISTS idx_ai_doc_review_run_files_tenant_id
ON ai_doc_review_run_files(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_doc_review_run_files_ai_doc_review_id
ON ai_doc_review_run_files(ai_doc_review_id);
CREATE INDEX IF NOT EXISTS idx_ai_doc_review_run_files_document_file_id
ON ai_doc_review_run_files(document_file_id);


END;
