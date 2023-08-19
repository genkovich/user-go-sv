-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
   id uuid PRIMARY KEY,
   login VARCHAR(255) NOT NULL,
   role VARCHAR(30) NOT NULL,
   password_hash VARCHAR(255) NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
