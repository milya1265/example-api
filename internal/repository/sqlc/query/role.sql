-- name: GetRole    :one
SELECT (role) FROM users_role WHERE user_id = $1;

-- name: AddRole :exec
INSERT INTO users_role (user_id, role) VALUES ($1, $2);

-- name: DeleteRole :exec
DELETE FROM users_role WHERE user_id = $1;