BEGIN;


ALTER TABLE IF EXISTS object ADD COLUMN run_version UUID;

UPDATE object AS o
SET run_version = sub.run_version
FROM (
    SELECT DISTINCT cluster_id, gen_random_uuid() AS run_version
    FROM object
) AS sub
WHERE o.cluster_id = sub.cluster_id;


ALTER TABLE object
ALTER COLUMN run_version SET NOT NULL;


CREATE INDEX IF NOT EXISTS idx_object_run_version ON object (run_version);


END;
