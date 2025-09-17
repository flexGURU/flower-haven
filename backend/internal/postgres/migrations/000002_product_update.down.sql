ALTER TABLE "products" DROP COLUMN "has_stems";
ALTER TABLE "products" DROP COLUMN "is_message_card";
ALTER TABLE "products" DROP COLUMN "is_flowers";
ALTER TABLE "products" DROP COLUMN "is_add_on";
ALTER TABLE "product_stems" DROP CONSTRAINT "product_stems_product_id_fkey";

DROP TABLE IF EXISTS product_stems;

ALTER TABLE "orders" DROP COLUMN "delivery_date";
ALTER TABLE "orders" DROP COLUMN "time_slot";
ALTER TABLE "orders" DROP COLUMN "by_admin";

ALTER TABLE "order_items" DROP COLUMN "stem_id";
ALTER TABLE "order_items" DROP CONSTRAINT "order_items_stem_id_fkey";
ALTER TABLE "order_items" DROP COLUMN "payment_method";
ALTER TABLE "order_items" DROP COLUMN "frequency";

ALTER TABLE "subscriptions" DROP COLUMN "by_admin";
ALTER TABLE "subscriptions" DROP COLUMN "stem_ids";
ALTER TABLE "subscriptions" DROP COLUMN "parent_order_id";
ALTER TABLE "subscriptions" DROP CONSTRAINT "subscriptions_parent_order_id_fkey";

ALTER TABLE "user_subscriptions" DROP COLUMN "frequency";
