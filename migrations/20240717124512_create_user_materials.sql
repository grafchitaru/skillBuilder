-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user_materials"
(
    user_id uuid NOT NULL REFERENCES users(id),
    material_id uuid NOT NULL REFERENCES materials(id),
    completed boolean NOT NULL DEFAULT false, -- выполнен ли материал
    PRIMARY KEY (user_id, material_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_materials;
-- +goose StatementEnd