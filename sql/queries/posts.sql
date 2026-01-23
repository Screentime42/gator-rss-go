-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
   $1,
   $2,
   $3,
   $4,
   $5,
   $6,
   $7
)
ON CONFLICT (url) DO NOTHING
RETURNING *;


-- name: GetPostsForUser :many
SELECT *
FROM posts p
JOIN feed_follows ff ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;