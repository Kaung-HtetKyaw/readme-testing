BEGIN;


CREATE TABLE IF NOT EXISTS password_auth (
    email TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    password_salt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


INSERT INTO password_auth (email, password_hash, password_salt, created_at, updated_at)
SELECT 
    email,
    password_hash,
    password_salt,
    created_at,
    updated_at
FROM user_account
ON CONFLICT (email) DO NOTHING;


ALTER TABLE user_account
DROP COLUMN password_hash,
DROP COLUMN password_salt;


CREATE TABLE IF NOT EXISTS oauth (
    email TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS email_verification_token (
    email TEXT PRIMARY KEY,
    value UUID NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


INSERT INTO email_verification_token (email, value, expiration, created_at, updated_at)
SELECT
  ua.email,
  vt.value,
  vt.expiration,
  vt.created_at,
  vt.updated_at
FROM user_account AS ua JOIN verify_token AS vt ON ua.id = vt.user_id;


DROP TABLE IF EXISTS verify_token;


END;
