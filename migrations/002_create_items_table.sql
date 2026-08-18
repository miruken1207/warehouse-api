-- +goose Up
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT
);

-- +goose Down
DROP TABLE IF EXISTS items;
