BEGIN;


ALTER TABLE personal_access_token RENAME TO credential;
ALTER TABLE credential
ADD CONSTRAINT uniq_credential_org_name UNIQUE (organization_id, name);
ALTER TABLE credential DROP CONSTRAINT uniq_personal_access_token_org_name;


ALTER TABLE credential ADD COLUMN type TEXT NOT NULL DEFAULT 'pat';
ALTER TABLE credential DROP COLUMN owner;
ALTER TABLE credential DROP COLUMN created_by;
ALTER TABLE credential DROP COLUMN updated_by;


ALTER TABLE repository ADD COLUMN credential_id UUID REFERENCES credential(id);
ALTER TABLE repository DROP COLUMN provider;
ALTER TABLE repository DROP COLUMN created_by;
ALTER TABLE repository DROP COLUMN updated_by;
ALTER TABLE repository DROP COLUMN namespace;


UPDATE repository AS r
SET credential_id = rc.personal_access_token_id
FROM repository_personal_access_token AS rc
WHERE rc.repository_id = r.id;


ALTER TABLE repository
ALTER COLUMN credential_id SET NOT NULL;


DROP TABLE IF EXISTS repository_personal_access_token;


END;
