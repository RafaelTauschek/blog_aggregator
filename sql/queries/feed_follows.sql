-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, NOW(), NOW(), $2, $3)
    RETURNING *
)
SELECT inserted_feed_follow.*, users.name AS user_name, feeds.name AS feed_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;


-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows;

-- name: GetFeedFollowForUser :many
SELECT feed_follows.*, users.name AS user_name, feeds.name AS feeds_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds on feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;


-- name: UnfollowFeed :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
