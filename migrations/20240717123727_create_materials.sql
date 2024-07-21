-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "materials"
(
    id uuid PRIMARY KEY NOT NULL,
    created_at timestamp(0) without time zone NOT NULL,
    updated_at timestamp(0) without time zone NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id),
    name text NOT NULL,
    description text,
    type text NOT NULL,
    xp integer NOT NULL,
    link text
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE materials;
-- +goose StatementEnd