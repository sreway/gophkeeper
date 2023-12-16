CREATE TABLE IF NOT EXISTS client(
    id UUID NOT NULL UNIQUE,
    active_profile UUID,
    PRIMARY KEY (id),
    FOREIGN KEY (active_profile) REFERENCES profiles(id)
    );