-- name: CreateTokens :exec
INSERT INTO jwt (user_id, access_token, refresh_token) VALUES ($1, $2, $3);

-- name: UpdateAccessToken :exec
UPDATE jwt SET access_token = $1 WHERE user_id = $2;

-- name: UpdateTokens :exec
UPDATE jwt SET access_token = $1, refresh_token = $2 WHERE user_id = $3;

-- name: GetTokens :one
SELECT * FROM jwt WHERE user_id = $1;

-- name: DeleteTokens :exec
DELETE FROM jwt WHERE user_id = $1;

