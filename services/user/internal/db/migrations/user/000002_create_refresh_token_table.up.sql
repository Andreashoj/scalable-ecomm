CREATE TABLE refresh_token (
    id UUID PRIMARY KEY,
    token varchar(255) DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users(id)
);
