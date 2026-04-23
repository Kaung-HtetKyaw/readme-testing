BEGIN;

CREATE TABLE IF NOT EXISTS payment_orders (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    service_application_id UUID NOT NULL REFERENCES service_applications(id) ON DELETE RESTRICT,
    provider TEXT NOT NULL, 
    amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0),
    currency TEXT NOT NULL, 
    status TEXT NOT NULL CHECK (status IN ('requires_payment', 'processing', 'paid', 'failed', 'cancelled', 'refunded', 'partially_refunded')),
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_payment_orders_tenant_id
ON payment_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_tenant_id_user_id
ON payment_orders(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_tenant_id_service_application_id
ON payment_orders(tenant_id, service_application_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_tenant_id_status_created_at
ON payment_orders(tenant_id, status, created_at);


CREATE TABLE IF NOT EXISTS payment_attempts (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    payment_order_id UUID NOT NULL REFERENCES payment_orders(id) ON DELETE CASCADE,
    provider_customer_id TEXT, 
    provider_payment_intent_id TEXT NOT NULL, 
    provider_payment_method_id TEXT, 
    provider_charge_id TEXT, 
    status TEXT NOT NULL CHECK (status IN ('requires_payment', 'processing', 'succeeded', 'failed', 'cancelled', 'refunded', 'partially_refunded')),
    failure_code TEXT,
    failure_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);


CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_payment_attempts_provider_payment_intent_id
ON payment_attempts(provider_payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_tenant_id
ON payment_attempts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_tenant_id_payment_order_id
ON payment_attempts(tenant_id, payment_order_id);
CREATE INDEX IF NOT EXISTS idx_payment_attempts_tenant_id_status_created_at
ON payment_attempts(tenant_id, status, created_at);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_payment_attempts_payment_order_open_attempt
ON payment_attempts(payment_order_id)
WHERE status IN ('requires_payment', 'processing');


CREATE TABLE IF NOT EXISTS payment_refunds (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    payment_attempt_id UUID NOT NULL REFERENCES payment_attempts(id) ON DELETE CASCADE,
    provider_refund_id TEXT NOT NULL, 
    amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0),
    status TEXT NOT NULL CHECK (status IN ('pending', 'succeeded', 'failed', 'cancelled')),
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    succeeded_at TIMESTAMPTZ
);


CREATE UNIQUE INDEX IF NOT EXISTS uniq_idx_payment_refunds_provider_refund_id
ON payment_refunds(provider_refund_id);
CREATE INDEX IF NOT EXISTS idx_payment_refunds_tenant_id
ON payment_refunds(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_refunds_tenant_id_payment_attempt_id
ON payment_refunds(tenant_id, payment_attempt_id);
CREATE INDEX IF NOT EXISTS idx_payment_refunds_tenant_id_status_created_at
ON payment_refunds(tenant_id, status, created_at);


CREATE TABLE IF NOT EXISTS payment_transactions (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    payment_order_id UUID NOT NULL REFERENCES payment_orders(id) ON DELETE CASCADE,
    payment_attempt_id UUID REFERENCES payment_attempts(id) ON DELETE SET NULL,
    provider TEXT NOT NULL,
    transaction_type TEXT NOT NULL CHECK (transaction_type IN ('authorization', 'capture', 'charge_succeeded', 'charge_failed', 'refund_created', 'refund_succeeded', 'refund_failed', 'chargeback')),
    amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0),
    currency TEXT NOT NULL,
    provider_charge_id TEXT, -- e.g. ch_...
    provider_refund_id TEXT, -- e.g. re_...
    provider_event_id TEXT, -- webhook event reference that caused this row
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_payment_transactions_tenant_id
ON payment_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_tenant_id_payment_order_id
ON payment_transactions(tenant_id, payment_order_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_tenant_id_payment_attempt_id
ON payment_transactions(tenant_id, payment_attempt_id);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_tenant_id_transaction_type_occurred_at
ON payment_transactions(tenant_id, transaction_type, occurred_at);
CREATE INDEX IF NOT EXISTS idx_payment_transactions_tenant_id_provider_event_id
ON payment_transactions(tenant_id, provider_event_id);


CREATE TABLE IF NOT EXISTS payment_provider_events (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE RESTRICT,
    provider TEXT NOT NULL,
    provider_event_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payment_order_id UUID REFERENCES payment_orders(id) ON DELETE SET NULL,
    payment_attempt_id UUID REFERENCES payment_attempts(id) ON DELETE SET NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, provider, provider_event_id)
);


CREATE INDEX IF NOT EXISTS idx_payment_provider_events_tenant_id
ON payment_provider_events(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payment_provider_events_tenant_id_event_type_created_at
ON payment_provider_events(tenant_id, event_type, created_at);
CREATE INDEX IF NOT EXISTS idx_payment_provider_events_tenant_id_payment_attempt_id
ON payment_provider_events(tenant_id, payment_attempt_id);


END;
