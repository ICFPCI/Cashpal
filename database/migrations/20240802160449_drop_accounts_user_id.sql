-- +goose Up
-- +goose StatementBegin
ALTER TABLE accounts DROP COLUMN user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE accounts ADD COLUMN user_id;
-- +goose StatementEnd
