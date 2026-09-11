-- name: CreateUser :one
INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?) RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;  

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;    

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;   

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;    

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?; 
