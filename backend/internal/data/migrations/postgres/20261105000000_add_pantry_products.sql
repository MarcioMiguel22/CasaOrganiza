-- +goose Up
CREATE TABLE IF NOT EXISTS "pantry_products" (
    "id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL,
    "updated_at" timestamptz NOT NULL,
    "name" character varying NOT NULL,
    "quantity" double precision NOT NULL DEFAULT 0,
    "min_quantity" double precision NOT NULL DEFAULT 0,
    "unit" character varying NOT NULL DEFAULT '',
    "expiry_date" character varying NOT NULL DEFAULT '',
    "purchase_price" double precision NOT NULL DEFAULT 0,
    "on_shopping_list" boolean NOT NULL DEFAULT false,
    "group_pantry_products" uuid NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "pantry_products_groups_pantry_products"
        FOREIGN KEY ("group_pantry_products") REFERENCES "groups" ("id")
        ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS "pantry_products_group_id" ON "pantry_products" ("group_pantry_products");
CREATE INDEX IF NOT EXISTS "pantry_products_expiry_date" ON "pantry_products" ("expiry_date");

-- +goose Down
DROP TABLE IF EXISTS "pantry_products";
