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

ALTER TABLE orders ADD COLUMN delivery_date TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE orders ADD COLUMN time_slot VARCHAR(50) NOT NULL DEFAULT 'anytime';
ALTER TABLE orders ADD COLUMN by_admin BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE order_items ADD COLUMN stem_id bigint NULL;
ALTER TABLE order_items ADD CONSTRAINT "order_items_stem_id_fkey" FOREIGN KEY ("stem_id") REFERENCES "product_stems" ("id");
ALTER TABLE order_items ADD COLUMN payment_method VARCHAR(50) NOT NULL DEFAULT 'one_time' CHECK (payment_method IN ('one_time', 'subscription'));
ALTER TABLE order_items ADD COLUMN frequency VARCHAR(50) NULL CHECK (frequency IN ('weekly', 'bi_weekly', 'monthly'));

ALTER TABLE subscriptions ADD COLUMN by_admin BOOLEAN NOT NULL DEFAULT true;
ALTER TABLE subscriptions ADD COLUMN stem_ids int[] NOT NULL DEFAULT '{}';
ALTER TABLE subscriptions ADD COLUMN parent_order_id bigint NULL;
ALTER TABLE subscriptions ADD CONSTRAINT "subscriptions_parent_order_id_fkey" FOREIGN KEY ("parent_order_id") REFERENCES "orders" ("id");


ALTER TABLE user_subscriptions ADD COLUMN frequency VARCHAR(50) NOT NULL DEFAULT 'monthly' CHECK (frequency IN ('weekly', 'bi_weekly', 'monthly'));