-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
                       id UUID PRIMARY KEY,
                       user_id UUID NOT NULL,
                       name TEXT NOT NULL,
                       description TEXT NOT NULL,
                       is_completed BOOLEAN NOT NULL,
                       created_at TIMESTAMP NOT NULL,
                       updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd

