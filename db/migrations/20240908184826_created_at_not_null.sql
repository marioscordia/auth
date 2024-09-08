-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ALTER COLUMN created_at
SET DEFAULT now();

ALTER TABLE users
ALTER COLUMN updated_at
SET DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
