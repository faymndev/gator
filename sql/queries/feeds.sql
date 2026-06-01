-- name: CreateFeed :one
insert into feeds (id, name, url, user_id)
values ($1, $2, $3, $4)
returning *;

-- name: GetFeeds :many
select feeds.*, users.name as user_name from feeds
left join users 
  on feeds.user_id = users.id;

-- name: GetFeedByUrl :one
select * from feeds
where feeds.url = $1;
