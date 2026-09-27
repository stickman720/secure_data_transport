CREATE TABLE IF NOT EXISTS users(
    id INTEGER    PRIMARY KEY  AUTOINCREMENT,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP    DEFAULT  CURRENT_TIMESTAMP,
    isactive      BOOLEAN      DEFAULT  false,
    role          VARCHAR(20)  NOT NULL 
);


CREATE INDEX lindex ON users (username , password_hash);



CREATE TABLE refreshtk (
    id INTEGER    PRIMARY KEY  AUTOINCREMENT,
    userid INTEGER    NOT NULL ,
    token VARCHAR(50)  NOT NULL UNIQUE,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX tiken_index ON refreshtk (token);