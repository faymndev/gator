-- name: CreateFeedFollow :one
with inserted_feed_follow as (
  insert into feed_follows (user_id, feed_id)
  values ($1, $2)
  returning *
) 
select inserted_feed_follow.*, users.name as user_name, feeds.name as feed_name
from inserted_feed_follow
inner join users on inserted_feed_follow.user_id = users.id
inner join feeds on inserted_feed_follow.feed_id = feeds.id;

-- name: GetFollowing :many
select feed_follows.user_id as user_id, feeds.* from feed_follows
inner join feeds on feed_follows.feed_id = feeds.id
where feed_follows.user_id = $1;
