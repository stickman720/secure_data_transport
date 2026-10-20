-- name: Login :one
SELECT * FROM users WHERE username = $1 AND password_hash = $2;

-- name: CreateRefreshToken :one
INSERT INTO refreshtk (userid,teamuser_id, token) VALUES (sqlc.arg("useris"),sqlc.arg("teamuser_id"),sqlc.arg("token"))
RETURNING *;

-- name: CleanRefreshToken :exec
DELETE FROM refreshtk
WHERE userid = (SELECT id FROM users WHERE username = $1);

-- name: ValidRefreshToken :one
SELECT u.id, u.role
FROM users u
JOIN refreshtk r ON r.userid = u.id
WHERE r.token = $1 AND u.isactive = TRUE AND r.createdat > NOW() - make_interval(hours => $2::int)
LIMIT 1;