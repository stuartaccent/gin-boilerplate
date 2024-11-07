-- name: CreateUser :one
-- create a new user
INSERT INTO auth_users (email, hashed_password, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
-- get a user by their email
SELECT *
FROM auth_users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
-- get a user by their id
SELECT *
FROM auth_users
WHERE id = $1
LIMIT 1;

-- name: SetUserPasswordByEmail :exec
-- set a user's password
UPDATE auth_users
SET hashed_password = $2
WHERE email = $1;