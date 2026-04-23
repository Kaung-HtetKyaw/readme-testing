BEGIN;


CREATE INDEX IF NOT EXISTS idx_object_cluster_id ON object (cluster_id);
CREATE INDEX IF NOT EXISTS idx_object_namespace ON object (namespace);
CREATE INDEX IF NOT EXISTS idx_object_kind ON object (kind);


END;
