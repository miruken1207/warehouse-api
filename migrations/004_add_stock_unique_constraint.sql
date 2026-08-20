-- +goose Up
ALTER TABLE stock ADD CONSTRAINT stock_warehouse_item_unique UNIQUE (warehouse_id, item_id);

-- +goose Down
ALTER TABLE stock DROP CONSTRAINT stock_warehouse_item_unique;
