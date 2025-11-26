-- name: GetPrize :one
SELECT * FROM prizes WHERE id = $1;

-- name: PrizesList :many
SELECT * FROM prizes ORDER BY id;

-- name: PrizesListPaginated :many
SELECT * FROM prizes ORDER BY id LIMIT $1 OFFSET $2;

-- name: CountPrizes :one
SELECT COUNT(*) FROM prizes;

-- name: PrizesByCategory :many
SELECT * FROM prizes WHERE category = $1 ORDER BY year DESC;

-- name: PrizesByYear :many
SELECT * FROM prizes WHERE year = $1 ORDER BY category;

-- name: GetPrizeWithLaureates :many
SELECT 
  p.id as prize_id, p.year, p.category,
  l.id as laureate_id, l.firstname, l.surname, l.motivation, l.share
FROM prizes p
LEFT JOIN prizes_to_laureates ptl ON p.id = ptl.prize_id
LEFT JOIN laureates l ON ptl.laureate_id = l.id
WHERE p.id = $1
ORDER BY l.id;

-- name: GetPrizesByCategoryWithLaureates :many
SELECT 
  p.id as prize_id, p.year, p.category,
  l.id as laureate_id, l.firstname, l.surname, l.motivation, l.share
FROM prizes p
LEFT JOIN prizes_to_laureates ptl ON p.id = ptl.prize_id
LEFT JOIN laureates l ON ptl.laureate_id = l.id
WHERE p.category = $1
ORDER BY p.year DESC, l.id;

-- name: AddPrize :batchone
INSERT INTO prizes (year, category)
VALUES ($1, $2)
RETURNING *;

-- name: AddPrizeSingle :one
INSERT INTO prizes (year, category)
VALUES ($1, $2)
RETURNING *;

-- name: UpdatePrize :one
UPDATE prizes
SET year = $2, category = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePrize :exec
DELETE FROM prizes WHERE id = $1;

-- name: GetCategories :many
SELECT DISTINCT category FROM prizes ORDER BY category;

-- name: GetLaureatesByPrizeId :many
SELECT l.* 
FROM laureates l
INNER JOIN prizes_to_laureates ptl ON l.id = ptl.laureate_id
WHERE ptl.prize_id = $1
ORDER BY l.id;