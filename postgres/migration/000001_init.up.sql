BEGIN;


CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS user_account (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    role_name TEXT NOT NULL,
    email TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    password_salt TEXT NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT false,
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS verify_token (
    user_id UUID PRIMARY KEY REFERENCES user_account(id),
    value UUID NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS forgot_password_token (
    user_id UUID PRIMARY KEY REFERENCES user_account(id),
    value UUID NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS invited_user (
    email TEXT PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    role_name TEXT NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS cluster_group (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_cluster_group_org_name UNIQUE (organization_id, name)
);


CREATE TABLE IF NOT EXISTS cluster (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    cluster_group_id UUID NOT NULL REFERENCES cluster_group(id)
);


CREATE TABLE IF NOT EXISTS cluster_info (
    cluster_id UUID PRIMARY KEY REFERENCES cluster(id),
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    platform TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS object (
    id UUID NOT NULL,
    cluster_id UUID NOT NULL REFERENCES cluster(id),
    namespace TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    resource_version TEXT NOT NULL,
    kind TEXT NOT NULL,
    raw JSONB,
    health_status TEXT NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cluster_id, id)
);


CREATE TABLE IF NOT EXISTS deprecated_api_group (
    cluster_id UUID PRIMARY KEY REFERENCES cluster(id),
    deprecated_apis JSONB
);


CREATE TABLE IF NOT EXISTS sandbox (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS sandbox_cluster_info (
    sandbox_id UUID PRIMARY KEY REFERENCES sandbox(id),
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    platform TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS sandbox_object (
    sandbox_id UUID REFERENCES sandbox(id),
    id UUID NOT NULL,
    namespace TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    resource_version TEXT NOT NULL,
    kind TEXT NOT NULL,
    raw JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (sandbox_id, id)
);


CREATE TABLE IF NOT EXISTS sandbox_deprecated_api_group (
    sandbox_id UUID PRIMARY KEY REFERENCES sandbox(id),
    deprecated_apis JSONB
);


CREATE TABLE IF NOT EXISTS cluster_tag (
    cluster_id UUID REFERENCES cluster(id),
    name TEXT,
    PRIMARY KEY (cluster_id, name)
);


CREATE TABLE IF NOT EXISTS cluster_token (
    value TEXT PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_cluster_token_org_name UNIQUE (organization_id, name)
);


CREATE TABLE IF NOT EXISTS credential (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    provider TEXT NOT NULL,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    encrypted_value TEXT NOT NULL,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uniq_credential_org_name UNIQUE (organization_id, name)
);


END;
