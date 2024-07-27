-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id serial primary key,
    username varchar(255) not null,
    email varchar(255) not null,
    role varchar(255) not null,
    password varchar(255) not null,
    status varchar(255) default 'active',
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
