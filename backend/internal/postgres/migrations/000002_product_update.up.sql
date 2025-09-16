ALTER TABLE products ADD COLUMN has_stems BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE products ADD COLUMN is_message_card BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE products ADD COLUMN is_flowers BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE products ADD COLUMN is_add_on BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE "product_stems" (
    "id" bigserial PRIMARY KEY,
    "product_id" bigint NOT NULL,
    "stem_count" bigint NOT NULL DEFAULT 0,
    "price" decimal(10, 2) NOT NULL DEFAULT 0,
    
    CONSTRAINT "product_stems_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
);

CREATE INDEX idx_product_stems_product_id ON product_stems (product_id);