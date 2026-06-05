-- name: CreatePost :one
insert into posts (title, url, description, published_at, feed_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetPosts :many
select * from posts
where feed_id = $1
order by published_at desc
limit $2;
