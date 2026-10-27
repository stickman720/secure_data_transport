-- name: CreateTeamUser :one
INSERT INTO teamuser (username , team_id , token) VALUES ($1 , (SELECT id FROM teams WHERE title = $2) , $3) RETURNING *;

-- name: GetTeamAndUser :one
SELECT 
    t.title,
    t.rate_limit_per_user,
    t.storage_space,
    t.transfer,
    tu.username,
    tu.id
FROM teams t 
JOIN teamuser tu ON teamuser.team_id = t.id
WHERE teamuser.token = $1;

-- name: DeleteTeamUser :execrows
DELETE FROM teamuser tu
USING teams t
WHERE t.id = tu.team_id AND 
tu.username = $1 AND t.userid = sqlc.arg(owner_id) ;



-- name: GetAllTeamUsers :many
SELECT 
    tu.username,
    tu.id
FROM teamuser tu
JOIN teams t ON tu.team_id = t.id 
JOIN users u ON u.id = t.userid
WHERE u.id = $1 AND t.title = $2;


-- name: LoginTeamUser :one
SELECT username , id from teamuser WHERE token = $1;




