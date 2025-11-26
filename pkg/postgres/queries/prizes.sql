-- name: PrizesList :many
SELECT
  prizes.*,
  COALESCE(
    jsonb_agg(to_jsonb(laureates) ORDER BY laureates.id) FILTER (WHERE laureates.id IS NOT NULL),
    '[]'
  ) AS laureates
FROM prizes
LEFT JOIN laureates ON laureates.id = prizes.laureate_id
GROUP BY prizes.id
ORDER BY prizes.id;

-- name: AddPrize :batchone
INSERT INTO prizes (year, category)
VALUES ($1, $2)
RETURNING *;