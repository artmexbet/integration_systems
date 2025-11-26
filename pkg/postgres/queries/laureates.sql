-- name: GetLaureate :one
SELECT * FROM laureates
         WHERE id = $1;

-- name: ListLaureates :many
SELECT * FROM laureates
            ORDER BY id;

-- name: ListLaureatesPaginated :many
SELECT * FROM laureates
            ORDER BY id
            LIMIT $1 OFFSET $2;

-- name: CountLaureates :one
SELECT COUNT(*) FROM laureates;

-- name: CreateLaureate :batchone
INSERT INTO laureates (id, firstname, surname, motivation, share)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateLaureateSingle :one
INSERT INTO laureates (id, firstname, surname, motivation, share)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateLaureate :one
UPDATE laureates
SET firstname = $2, surname = $3, motivation = $4, share = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteLaureate :exec
DELETE FROM laureates WHERE id = $1;

-- name: LinkLaureateToPrize :batchexec
INSERT INTO prizes_to_laureates (prize_id, laureate_id)
VALUES ($1, $2);

-- name: LinkLaureateToPrizeSingle :exec
INSERT INTO prizes_to_laureates (prize_id, laureate_id)
VALUES ($1, $2);