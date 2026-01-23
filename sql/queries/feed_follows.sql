-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
   INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
   VALUES (
      $1,
      $2,
      $3,
      $4
   )
   RETURNING *
)

SELECT
   inserted_feed_follow.*,
   users.name AS user_name,
   feeds.name AS feed_name
   
FROM inserted_feed_follow
JOIN users ON users.id = inserted_feed_follow.user_id
JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;


-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;


-- name: GetFeedFollowsForUser :many
SELECT 
   f.name AS feed_name, 
   f.url AS feed_url, 
   u.name AS user_name   
FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1;


-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE user_id = $1
   AND feed_id = $2;