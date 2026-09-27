CREATE TABLE IF NOT EXISTS users (
    id            SERIAL       PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    isactive      BOOLEAN      DEFAULT FALSE,
    role          VARCHAR(20)  NOT NULL
);




CREATE TABLE refreshtk (
    id     SERIAL       PRIMARY KEY,
    userid INTEGER      NOT NULL,
    token  VARCHAR(50)  NOT NULL UNIQUE,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
);
