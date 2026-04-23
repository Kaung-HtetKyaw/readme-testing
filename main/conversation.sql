BEGIN;

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    provider TEXT NOT NULL,
    provider_conversation_id TEXT NOT NULL,
    provider_contact_id TEXT NOT NULL,
    channel_id TEXT,
    channel_name TEXT,
    channel_source TEXT,
    channel_meta JSONB,
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    provider_assignee_id TEXT,
    status TEXT, -- open, pending, resolved, closed or statuses from "conversation_status_changed" or any event submitted by the provider
    lifecycle TEXT NOT NULL CHECK (lifecycle IN ('open', 'pending', 'resolved', 'closed')),
    contact_status TEXT,
    opened_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, provider, provider_conversation_id)
);


CREATE INDEX IF NOT EXISTS idx_conversations_tenant_id
ON conversations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_conversations_provider_contact_id
ON conversations(provider_contact_id);
CREATE INDEX IF NOT EXISTS idx_conversations_assignee_id
ON conversations(assignee_id);
CREATE INDEX IF NOT EXISTS idx_conversations_lifecycle
ON conversations(lifecycle);
CREATE INDEX IF NOT EXISTS idx_conversations_channel_source
ON conversations(channel_source);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at
ON conversations(created_at);


CREATE TABLE IF NOT EXISTS staff_provider_agents (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    provider TEXT NOT NULL,
    provider_agent_id TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, provider, provider_agent_id)
);


CREATE INDEX IF NOT EXISTS idx_staff_provider_agents_tenant_id
ON staff_provider_agents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_staff_provider_agents_user_id
ON staff_provider_agents(user_id);
CREATE INDEX IF NOT EXISTS idx_staff_provider_agents_provider
ON staff_provider_agents(provider);
CREATE INDEX IF NOT EXISTS idx_staff_provider_agents_is_active
ON staff_provider_agents(is_active);


CREATE TABLE IF NOT EXISTS conversation_service_applications (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    service_application_id UUID NOT NULL REFERENCES service_applications(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id, service_application_id)
);


CREATE INDEX IF NOT EXISTS idx_conversation_service_applications_tenant_id
ON conversation_service_applications(tenant_id);
CREATE INDEX IF NOT EXISTS idx_conversation_service_applications_conversation_id
ON conversation_service_applications(conversation_id);
CREATE INDEX IF NOT EXISTS idx_conversation_service_applications_service_application_id
ON conversation_service_applications(service_application_id);


CREATE TABLE IF NOT EXISTS conversation_crm_profiles (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    client_location TEXT,
    client_nationality TEXT,
    current_visa TEXT,
    dtv_purpose TEXT,
    dtv_package TEXT,
    submission_country TEXT,
    client_interested_in TEXT,
    client_urgency TEXT,
    client_buying_segment TEXT,
    client_buying_intent TEXT,
    current_step TEXT,
    source TEXT NOT NULL CHECK (source IN ('agent', 'ai', 'system', 'webhook')),
    updated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (conversation_id)
);


CREATE INDEX IF NOT EXISTS idx_conversation_crm_profiles_tenant_id
ON conversation_crm_profiles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profiles_client_interested_in
ON conversation_crm_profiles(client_interested_in);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profiles_client_urgency
ON conversation_crm_profiles(client_urgency);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profiles_client_buying_segment
ON conversation_crm_profiles(client_buying_segment);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profiles_client_buying_intent
ON conversation_crm_profiles(client_buying_intent);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profiles_submission_country
ON conversation_crm_profiles(submission_country);


CREATE TABLE IF NOT EXISTS conversation_crm_profile_history (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    field_name TEXT NOT NULL,
    value TEXT,
    source TEXT NOT NULL CHECK (source IN ('agent', 'ai', 'system', 'webhook')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL
);


CREATE INDEX IF NOT EXISTS idx_conversation_crm_profile_history_tenant_id
ON conversation_crm_profile_history(tenant_id);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profile_history_conversation_id
ON conversation_crm_profile_history(conversation_id);
CREATE INDEX IF NOT EXISTS idx_conversation_crm_profile_history_conversation_field_created_at
ON conversation_crm_profile_history(conversation_id, field_name, created_at);


END;
