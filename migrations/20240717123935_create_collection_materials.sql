-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "collection_materials"
(
    collection_id uuid NOT NULL REFERENCES collections(id),
    material_id uuid NOT NULL REFERENCES materials(id),
    PRIMARY KEY (collection_id, material_id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE collection_materials;
-- +goose StatementEnd