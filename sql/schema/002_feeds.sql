-- +goose Up
create table feeds (
  id uuid primary key,
  created_at timestamp not null default now (),
  updated_at timestamp not null default now (),
  name text unique not null,
  url text unique not null,
  user_id uuid not null,
  constraint fk_user foreign key (user_id) references users (id) on delete cascade
);

-- +goose Down
drop table feeds;
