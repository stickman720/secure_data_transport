CREATE TABLE files (
    id SERIAL PRIMARY KEY ,
    team_id INTEGER,
    filename VARCHAR(250),
    size INT,
    type VARCHAR(50),
    address VARCHAR(100),
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

