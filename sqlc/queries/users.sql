-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    id,
    "fullName",
    "phoneNumber",
    "passwordHash"
)
VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPhoneNumber :one
SELECT * FROM users WHERE "phoneNumber" = $1;

-- name: UpdateUser :one
UPDATE users SET
    "fullName" = COALESCE($2, "fullName"),
    "phoneNumber" = COALESCE($3, "phoneNumber"),
    "passwordHash" = COALESCE($4, "passwordHash"),
    "updatedAt" = now()
WHERE id = $1
RETURNING *;