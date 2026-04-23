BEGIN;


CREATE TABLE IF NOT EXISTS cluster_upgrade (
    cluster_id UUID PRIMARY KEY REFERENCES cluster(id) ON DELETE CASCADE,
    next_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS upgradeable_component (
    cluster_id UUID NOT NULL REFERENCES cluster(id) ON DELETE CASCADE,
    object_id UUID NOT NULL,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    next_compatible BOOLEAN NOT NULL,
    min_compatible_version TEXT NOT NULL,
    max_compatible_version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (cluster_id, object_id) REFERENCES object(cluster_id, id) ON DELETE CASCADE,
    PRIMARY KEY (cluster_id, object_id, name)
);


CREATE TABLE IF NOT EXISTS unmatched_component (
    cluster_id UUID NOT NULL REFERENCES cluster(id) ON DELETE CASCADE,
    object_id UUID NOT NULL,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (cluster_id, object_id) REFERENCES object(cluster_id, id) ON DELETE CASCADE,
    PRIMARY KEY(cluster_id, object_id, name)
);


END;
