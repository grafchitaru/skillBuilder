-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "type_materials"
(
    id uuid PRIMARY KEY NOT NULL,
    created_at timestamp(0) without time zone NOT NULL,
    updated_at timestamp(0) without time zone NOT NULL,
    name text NOT NULL,
    characteristic text NOT NULL,
    xp integer NOT NULL
);
INSERT INTO "type_materials" (id, created_at, updated_at, name, characteristic, xp)
VALUES
    ('1ef49c5e-fc3e-6b7e-9532-53fb33479b19', NOW(), NOW(), 'книга', 'страница', 1),
    ('1ef49c5f-643c-6226-8913-f57081f12b8e', NOW(), NOW(), 'аудио-книга', 'час', 10),
    ('1ef49c5f-9f02-6680-abd1-41b757f22f2c', NOW(), NOW(), 'статья', 'штука', 3),
    ('1ef49c5f-d824-6f1c-b0d2-bb6c3997fabe', NOW(), NOW(), 'курс', 'урок', 10),
    ('1ef49c60-0fd2-6a36-85d5-e1237e109465', NOW(), NOW(), 'видеоролик', 'час', 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE type_materials;
-- +goose StatementEnd
