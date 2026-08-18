-- +goose Up
CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    location TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS warehouses;
