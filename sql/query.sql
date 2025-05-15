-- name: GetUser :one
SELECT * FROM "user" 
WHERE userId = $1 LIMIT 1;

-- name: ListUsers :many
SELECT userId, first_name, last_name, email, phone, age, "status" FROM "user"
ORDER BY first_name;

-- name: CreateUser :exec
INSERT INTO "user" (
    first_name, last_name, email, phone, age, "status"
    )values (
        $1, $2, $3, $4, $5, $6);

-- name: UpdateUser :exec
UPDATE "user" 
SET first_name = $1, last_name = $2, email = $3, phone = $4, age = $5, "status" = $6 
WHERE userId = $7;

-- name: DeleteUser :exec
DELETE FROM "user" 
WHERE userId = $1;
