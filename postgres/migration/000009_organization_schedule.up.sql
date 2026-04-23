BEGIN;


CREATE TABLE IF NOT EXISTS organization_schedule (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organization(id),
    name TEXT NOT NULL,
    payload JSONB,
    scheduled_for TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


END;
