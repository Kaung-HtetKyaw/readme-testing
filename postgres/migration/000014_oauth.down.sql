BEGIN;


ALTER TABLE user_account
ADD COLUMN password_hash TEXT,
ADD COLUMN password_salt TEXT;


UPDATE user_account AS ua
SET 
    password_hash = b.password_hash,
    password_salt = b.password_salt
FROM password_auth b
WHERE ua.email = b.email;


ALTER TABLE user_account
ALTER COLUMN password_hash SET NOT NULL,
ALTER COLUMN password_salt SET NOT NULL;


DROP TABLE IF EXISTS password_auth;
DROP TABLE IF EXISTS oauth;


CREATE TABLE IF NOT EXISTS verify_token (
    user_id UUID PRIMARY KEY REFERENCES user_account(id),
    value UUID NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


INSERT INTO verify_token (user_id, value, expiration, created_at, updated_at)
SELECT
  ua.id,
  evt.value,
  evt.expiration,
  evt.created_at,
  evt.updated_at
FROM user_account AS ua JOIN email_verification_token AS evt ON ua.email = evt.email;


DROP TABLE IF EXISTS  email_verification_token;


END;
