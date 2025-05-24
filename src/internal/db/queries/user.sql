-- name: GetUserByID :one
SELECT * FROM Users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM Users
ORDER BY name;

-- name: CreateUser :execresult
INSERT INTO Users (
  id, name, email, password
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateUser :execresult
UPDATE Users
SET name = ?, email = ?, password = ?
WHERE id = ?;

-- name: DeleteUser :execresult
DELETE FROM Users
WHERE id = ?;
