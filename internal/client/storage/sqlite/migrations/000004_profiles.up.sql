CREATE TABLE IF NOT EXISTS profiles (
    id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL UNIQUE,
    seal_id UUID NOT NULL UNIQUE,
    session_id UUID,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (seal_id) REFERENCES seals(id),
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);