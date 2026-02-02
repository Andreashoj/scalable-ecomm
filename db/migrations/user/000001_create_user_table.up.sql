CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    refresh_token_expires_at TIMESTAMP NOT NULL
);