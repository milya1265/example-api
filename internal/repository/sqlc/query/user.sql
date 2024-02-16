-- name: CreateUser :exec
INSERT INTO users (id, login, password_hash) VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT * FROM users WHERE login = $1;

-- name: UpdateUser :exec
UPDATE users SET login = $1, password_hash = $2 WHERE id = $3;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
