-- name: GetLastUpdate :one
SELECT 
  COALESCE(
    GREATEST(
      (SELECT MAX(updated_at) FROM laureates),
      (SELECT MAX(updated_at) FROM prizes)
    ),
    NOW()
  ) as last_update;

-- name: GetStats :one
SELECT 
  (SELECT COUNT(*) FROM laureates) as laureates_count,
  (SELECT COUNT(*) FROM prizes) as prizes_count,
  (SELECT COUNT(DISTINCT category) FROM prizes) as categories_count;
