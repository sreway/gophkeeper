CREATE TABLE IF NOT EXISTS sessions(
    id UUID NOT NULL UNIQUE,
    encrypted_token  BYTEA NOT NULL,
    PRIMARY KEY (id)
    );