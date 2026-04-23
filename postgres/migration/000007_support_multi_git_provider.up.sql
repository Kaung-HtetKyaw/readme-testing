BEGIN;


ALTER TABLE credential RENAME TO personal_access_token;
ALTER TABLE personal_access_token
ADD CONSTRAINT uniq_personal_access_token_org_name UNIQUE (organization_id, provider, name);
ALTER TABLE personal_access_token DROP CONSTRAINT uniq_credential_org_name;

    
ALTER TABLE personal_access_token DROP COLUMN type;
ALTER TABLE personal_access_token 
ADD COLUMN owner TEXT NOT NULL DEFAULT 'unknown';
ALTER TABLE personal_access_token 
ADD COLUMN created_by UUID REFERENCES user_account(id);
ALTER TABLE personal_access_token 
ADD COLUMN updated_by UUID REFERENCES user_account(id);


UPDATE personal_access_token AS p
SET created_by = u.id
FROM user_account AS u
WHERE u.organization_id = p.organization_id 
AND u.role_name = 'owner';


UPDATE personal_access_token AS p
SET updated_by = u.id
FROM user_account AS u
WHERE u.organization_id = p.organization_id 
AND u.role_name = 'owner';


CREATE TABLE IF NOT EXISTS repository_personal_access_token (
    repository_id UUID REFERENCES repository(id),
    personal_access_token_id UUID REFERENCES personal_access_token(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (repository_id, personal_access_token_id)
);


INSERT INTO repository_personal_access_token (repository_id, personal_access_token_id)
SELECT id, credential_id FROM repository WHERE credential_id IS NOT NULL;


ALTER TABLE repository DROP COLUMN credential_id;
ALTER TABLE repository ADD COLUMN provider TEXT NOT NULL default 'github';
ALTER TABLE repository ADD COLUMN created_by UUID REFERENCES user_account(id);
ALTER TABLE repository ADD COLUMN updated_by UUID REFERENCES user_account(id);
ALTER TABLE repository ADD COLUMN namespace TEXT NOT NULL default '';


UPDATE repository AS r
SET created_by = u.id
FROM user_account AS u
WHERE u.organization_id = r.organization_id 
AND (u.role_name = 'owner' OR u.role_name = 'admin');


UPDATE repository AS r
SET updated_by = u.id
FROM user_account AS u
WHERE u.organization_id = r.organization_id 
AND (u.role_name = 'owner' OR u.role_name = 'admin');


END;
