-- name: GetLaureate :one
SELECT * FROM laureates
         WHERE id = $1;

-- name: ListLaureates :many
SELECT * FROM laureates
            ORDER BY id;

-- name: CreateLaureate :batchone
INSERT INTO laureates (id, firstname, surname, motivation, share)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: LinkLaureateToPrize :batchexec
INSERT INTO prizes_to_laureates (prize_id, laureate_id)
VALUES ($1, $2);