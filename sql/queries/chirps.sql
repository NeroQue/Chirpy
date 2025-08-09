-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2
)
RETURNING *;
--

-- name: GetAllChirps :many
SELECT * FROM chirps
order by created_at;
--

-- name: GetChirp :one
Select * from chirps
where id = $1;
--

-- name: DeleteChirp :one
DELETE
FROM chirps
WHERE id = $1
RETURNING *;
--