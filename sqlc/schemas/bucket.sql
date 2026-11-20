CREATE TABLE IF NOT EXISTS buckets (
    id SERIAL PRIMARY KEY ,
    name TEXT ,
    location TEXT ,
    team_id INTEGER ,
    FOREIGN KEY (team_id) REFERENCES teams(id)
)