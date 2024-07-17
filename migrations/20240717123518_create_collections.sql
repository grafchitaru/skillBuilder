-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "collections"
(
    id uuid PRIMARY KEY NOT NULL,
    created_at timestamp(0) without time zone NOT NULL,
    updated_at timestamp(0) without time zone NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id),
    name text NOT NULL,
    description text
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE collections;
-- +goose StatementEnd