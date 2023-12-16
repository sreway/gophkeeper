CREATE TABLE IF NOT EXISTS seals(
    id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL UNIQUE,
    encrypted_shares BYTEA NOT NULL,
    recovery_share BYTEA NOT NULL,
    required_shares INTEGER NOT NULL,
    total_shares INTEGER NOT NULL,
    hash_master_password  VARCHAR(255) NOT NULL,
    hash_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    PRIMARY KEY (id)
);