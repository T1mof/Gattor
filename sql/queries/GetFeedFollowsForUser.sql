-- name: GetFeedFollowsForUser :many
SELECT 
    feeds.*,
    users.name as user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;
