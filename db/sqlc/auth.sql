-- name: CreateUser :one
-- create a new user
INSERT INTO auth_users (email, hashed_password, first_name, last_name)
VALUES (@email, @hashed_password, @first_name, @last_name)
RETURNING *;

-- name: GetUserByEmail :one
-- get a user by their email
SELECT *
FROM auth_users
WHERE email = @email
LIMIT 1;

-- name: GetUserByID :one
-- get a user by their id
SELECT *
FROM auth_users
WHERE id = @id
LIMIT 1;
