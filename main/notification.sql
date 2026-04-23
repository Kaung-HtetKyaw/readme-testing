BEGIN;

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    code TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('queued', 'processing', 'sent', 'failed', 'cancelled')),
    resource_type TEXT, -- e.g. service_application, document_submission, conversation
    resource_id UUID,
    scheduled_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id
ON notifications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id_status_scheduled_at
ON notifications(tenant_id, status, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id_user_id_created_at
ON notifications(tenant_id, user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id_code_created_at
ON notifications(tenant_id, code, created_at);
CREATE INDEX IF NOT EXISTS idx_notifications_tenant_id_resource_type_resource_id
ON notifications(tenant_id, resource_type, resource_id);


CREATE TABLE IF NOT EXISTS notification_targets (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel TEXT NOT NULL CHECK (channel IN ('email', 'push', 'sms', 'in_app')),
    target_value TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, user_id, channel, target_value)
);


CREATE INDEX IF NOT EXISTS idx_notification_targets_tenant_id
ON notification_targets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notification_targets_tenant_id_user_id_channel_is_active
ON notification_targets(tenant_id, user_id, channel, is_active);


CREATE TABLE IF NOT EXISTS notification_preferences (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel TEXT NOT NULL CHECK (channel IN ('email', 'push', 'sms', 'in_app')),
    notification_code TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, user_id, channel, notification_code)
);


CREATE INDEX IF NOT EXISTS idx_notification_preferences_tenant_id
ON notification_preferences(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notification_preferences_tenant_id_user_id_channel
ON notification_preferences(tenant_id, user_id, channel);


CREATE TABLE IF NOT EXISTS notification_deliveries (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    channel TEXT NOT NULL CHECK (channel IN ('email', 'push', 'sms', 'in_app')),
    target_id UUID REFERENCES notification_targets(id) ON DELETE SET NULL,
    attempt_no INTEGER NOT NULL DEFAULT 1,
    provider_message_id TEXT,
    delivery_status TEXT NOT NULL CHECK (delivery_status IN ('queued', 'sent', 'delivered', 'failed', 'bounced')),
    error_message TEXT,
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_notification_deliveries_tenant_id
ON notification_deliveries(tenant_id);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_tenant_id_notification_id
ON notification_deliveries(tenant_id, notification_id);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_tenant_id_channel_delivery_status_created_at
ON notification_deliveries(tenant_id, channel, delivery_status, created_at);
CREATE INDEX IF NOT EXISTS idx_notification_deliveries_tenant_id_provider_message_id
ON notification_deliveries(tenant_id, provider_message_id);


END;
