-- name: CreateUser :one
insert into users (id, name)
values ($1, $2)
returning *;

-- name: GetUser :one
select * from users
where users.name = $1
limit 1;

-- name: GetUsers :many
select * from users;

