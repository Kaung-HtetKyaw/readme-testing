BEGIN;


ALTER TABLE cluster 
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN last_observed;


ALTER TABLE cluster_info
DROP COLUMN created_at,
DROP COLUMN updated_at;


ALTER TABLE deprecated_api_group
DROP COLUMN created_at,
DROP COLUMN updated_at;


END;
