-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user_collections"
(
    user_id uuid NOT NULL REFERENCES users(id),
    collection_id uuid NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, collection_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_collections;
-- +goose StatementEnd