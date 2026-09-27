-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $1 WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2) AS user_exists;

-- name: VerifyUser :one
UPDATE users
SET isactive = TRUE
WHERE id = $1 AND isactive = FALSE
RETURNING *;

-- name: Login :one
SELECT * FROM users WHERE username = $1 AND password_hash = $2;

-- name: CreateRefreshToken :one
INSERT INTO refreshtk (userid, token) VALUES ($1, $2)
RETURNING *;

-- name: CleanRefreshToken :exec
DELETE FROM refreshtk
WHERE userid = (SELECT id FROM users WHERE username = $1);

-- name: ValidRefreshToken :one
SELECT u.username, u.role
FROM users u
JOIN refreshtk r ON r.userid = u.id
WHERE r.token = $1 AND u.isactive = TRUE
LIMIT 1;