BEGIN;


DROP INDEX IF EXISTS idx_object_run_version;


ALTER TABLE IF EXISTS object DROP COLUMN run_version;
  

END;
