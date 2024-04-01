-- name: GetList :many
SELECT * FROM List WHERE userid = ?;

-- name: PostList :one
INSERT INTO List (
  todo,
  userid
) VALUES (
  ?,
  ?
)
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM List WHERE id = ? AND userid = ?;

-- name: LoginUser :one
SELECT * FROM User WHERE email = ? AND password = ?;

-- name: CreateUser :one
INSERT INTO User (
  email,
  password
) VALUES (
  ?,
  ?
)
RETURNING *;
