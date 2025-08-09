-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    DEFAULT
)
RETURNING *;
--

-- name: DeleteAllUsers :exec
DELETE FROM users;
--

-- name: GetUserByEmail :one
SELECT * from users where email = $1;
--

-- name: UpdateUser :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = now()
WHERE id = $3
RETURNING *;
--

-- name: UpgradeUser :one
UPDATE users
set is_chirpy_red = true, updated_at = now()
where id = $1
RETURNING *;