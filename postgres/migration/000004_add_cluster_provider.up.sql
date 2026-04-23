BEGIN;


ALTER TABLE cluster_info ADD COLUMN provider TEXT NOT NULL DEFAULT 'unknown';


END;
