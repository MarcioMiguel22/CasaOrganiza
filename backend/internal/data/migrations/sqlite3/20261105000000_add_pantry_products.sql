-- +goose Up
CREATE TABLE IF NOT EXISTS pantry_products (
    id                     uuid     NOT NULL PRIMARY KEY,
    created_at             datetime NOT NULL,
    updated_at             datetime NOT NULL,
    name                   text     NOT NULL,
    quantity               real     NOT NULL DEFAULT 0,
    min_quantity           real     NOT NULL DEFAULT 0,
    unit                   text     NOT NULL DEFAULT '',
    expiry_date            text     NOT NULL DEFAULT '',
    purchase_price         real     NOT NULL DEFAULT 0,
    on_shopping_list       bool     NOT NULL DEFAULT false,
    group_pantry_products  uuid     NOT NULL
        CONSTRAINT pantry_products_groups_pantry_products
        REFERENCES groups ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS pantry_products_group_id ON pantry_products (group_pantry_products);
CREATE INDEX IF NOT EXISTS pantry_products_expiry_date ON pantry_products (expiry_date);

-- +goose Down
DROP TABLE IF EXISTS pantry_products;
