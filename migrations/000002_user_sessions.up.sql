CREATE TABLE graphics.user_sessions(
    token_hash BYTEA PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES graphics.users (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);