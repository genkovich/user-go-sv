-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD CONSTRAINT idx_unique_login UNIQUE (login); ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT idx_unique_login;
-- +goose StatementEnd
