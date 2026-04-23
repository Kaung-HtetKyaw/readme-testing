BEGIN;

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    actor_type TEXT NOT NULL CHECK (actor_type IN ('user', 'service_account', 'scheduled_job', 'system')),
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    target_type TEXT NOT NULL,
    target_id TEXT NOT NULL,
    correlation_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id
ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id_created_at
ON audit_logs(tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id_target_type_target_id_created_at
ON audit_logs(tenant_id, target_type, target_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id_action_created_at
ON audit_logs(tenant_id, action, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant_id_actor_user_id_created_at
ON audit_logs(tenant_id, actor_user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_correlation_id
ON audit_logs(correlation_id);


END;
