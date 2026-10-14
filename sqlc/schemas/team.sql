CREATE TABLE  teams (
    id                  SERIAL       PRIMARY KEY,
    title               VARCHAR(50)  NOT NULL UNIQUE,
    user_count          INTEGER      NOT NULL,
    rate_limit_per_user INTEGER      NOT NULL,
    storage_space       INTEGER      NOT NULL,
    transfer          INTEGER      NOT NULL,
    userid              SERIAL       NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
);