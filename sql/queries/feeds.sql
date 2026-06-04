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

-- name: MarkFeedFetched :exec
update feeds
set last_fetched_at = $1, updated_at = $2 
where feeds.id = $3;

-- name: GetNextFeedToFetch :one
select * from feeds
order by feeds.last_fetched_at asc nulls first
limit 1;
