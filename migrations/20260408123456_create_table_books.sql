-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS books (
    id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title   TEXT NOT NULL,
    author  TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at  TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd