-- name: CreateTeam :one
INSERT INTO teams (title , user_count ,rate_limit_per_user , storage_space , treansfear , userid) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *;

-- name: GetTeams :many
SELECT * FROM teams WHERE userid = $1 ;



