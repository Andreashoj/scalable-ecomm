CREATE TABLE access_token (
   id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    refresh_token_id UUID REFERENCES users (id)
);
