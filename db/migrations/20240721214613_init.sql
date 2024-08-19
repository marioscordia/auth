-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id serial primary key,
    username varchar(255) not null unique,
    email varchar(255) not null unique,
    role varchar(255) not null,
    password varchar(255) not null,
    status varchar(255) default 'active',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
