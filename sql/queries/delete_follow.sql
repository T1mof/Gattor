-- name: DeleteFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
  AND feed_id IN (
      SELECT id FROM feeds WHERE url = $2
);