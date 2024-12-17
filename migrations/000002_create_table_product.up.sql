-- 000002_create_table_product.up.sql
CREATE TABLE IF NOT EXISTS products
(
    id           UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    product_name VARCHAR(255)   NOT NULL,
    price        DECIMAL(15, 2) NOT NULL CHECK (price > 0),
    quantity     INTEGER        NOT NULL CHECK (quantity >= 0),
    brand_id     UUID           NOT NULL REFERENCES brands (id),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_products_deleted_at ON products (deleted_at);
CREATE INDEX idx_products_id_brand ON products (id);