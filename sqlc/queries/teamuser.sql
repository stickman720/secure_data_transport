-- name: CreateTeamUser :one
INSERT INTO teamuser (username , team_id , token) VALUES ($1 , (SELECT id FROM teams WHERE title = $2) , $3) RETURNING *;

-- name: GetTeamAndUser :one
SELECT 
    t.title,
    t.rate_limit_per_user,
    t.storage_space,
    t.transfer,
    t.userid,
    tu.username,
    tu.id
FROM teams t 
JOIN teamuser tu ON teamuser.team_id = t.id
WHERE teamuser.token = $1;

-- name: DeleteTeamUser :exec
DELETE FROM teamuser WHERE username = $1;