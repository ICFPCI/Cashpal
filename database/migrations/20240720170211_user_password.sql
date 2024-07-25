-- +goose Up
-- +goose StatementBegin
ALTER TABLE Users
RENAME COLUMN password_hash TO password;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE Users
RENAME COLUMN password TO password_hash;
-- +goose StatementEnd
