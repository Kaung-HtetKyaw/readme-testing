BEGIN;


ALTER TABLE upgradeable_component
DROP COLUMN run_version;

ALTER TABLE unmatched_component
DROP COLUMN run_version;


END;
