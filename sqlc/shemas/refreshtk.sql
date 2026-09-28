CREATE TABLE refreshtk (
    id     SERIAL       PRIMARY KEY,
    userid INTEGER      NOT NULL,
    token  VARCHAR(50)  NOT NULL UNIQUE,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
);