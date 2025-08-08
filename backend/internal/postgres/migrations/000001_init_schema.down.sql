ALTER TABLE "user_subscriptions" DROP CONSTRAINT "user_subscription_user_id_fkey";
ALTER TABLE "user_subscriptions" DROP CONSTRAINT "user_subscription_subscription_id_fkey";
ALTER TABLE "subscription_deliveries" DROP CONSTRAINT "subscription_deliveries_user_subscription_id_fkey";
ALTER TABLE "products" DROP CONSTRAINT "products_category_id_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT "order_item_order_id_fkey";
ALTER TABLE "order_items" DROP CONSTRAINT "order_item_product_id_fkey";
ALTER TABLE "payments" DROP CONSTRAINT "payment_order_id_fkey";

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS user_subscriptions;
DROP TABLE IF EXISTS subscription_deliveries;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS payments;