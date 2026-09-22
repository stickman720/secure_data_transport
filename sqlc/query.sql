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


-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?) AS user_exists;

-- name: VerifyUser :one
update users set isactive = 1 where id = ? and isactive = 0 returning *;