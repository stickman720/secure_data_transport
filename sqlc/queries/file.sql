-- name: NewFile :one 
INSERT INTO files (team_id , filename , size , type , address) VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- name: GetFiles :many 
SELECT 
    id,
    filename,
    size,
    type
FROM files
WHERE team_id = $1;


-- name: GetFileNameAddress :one
SELECT 
    filename,
    address 
FROM files
WHERE id = $1;


-- name: AddressExists :one 
SELECT EXISTS(SELECT 1 FROM files WHERE address = $1) AS address_exists ; 



-- name: DeleteFile :execrows
DELETE FROM files WHERE id = $1;



-- name: UpdateFileSize :execrows
UPDATE files SET size = $1 WHERE id = $2;




