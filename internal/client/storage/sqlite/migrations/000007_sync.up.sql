CREATE TABLE IF NOT EXISTS sync(
    user_id UUID NOT NULL,
    last_sync TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);