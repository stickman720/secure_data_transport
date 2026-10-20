CREATE TABLE refreshtk (
    id     SERIAL       PRIMARY KEY,
    userid INTEGER ,
    token  VARCHAR(50)  NOT NULL UNIQUE,
    createdat TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    teamuser_id SERIAL ,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (teamuser_id) REFERENCES teamuser(id) ON DELETE CASCADE,
    CONSTRAINT refereshtk_owner_chack CHECK ((userid IS NOT NULL AND teamuser_id IS NULL) OR (userid IS NULL AND teamuser_id IS NOT NULL))
);