CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "email" varchar(50) UNIQUE NOT NULL,
  "address" varchar(255) NULL,
  "phone_number" varchar(50) UNIQUE NOT NULL,
  "refresh_token" text NULL,
  "password" varchar(255) NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  "is_active" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "subscriptions" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "description" text NOT NULL,
  "product_ids" int[] NOT NULL DEFAULT '{}',
  "add_ons" int[] NOT NULL DEFAULT '{}',
  "price" decimal(10,2) NOT NULL DEFAULT 0,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_subscriptions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "subscription_id" bigint NOT NULL,
  "day_of_week" SMALLINT NOT NULL,
  "status" boolean NOT NULL DEFAULT true,
  "start_date" TIMESTAMPTZ NOT NULL,
  "end_date" TIMESTAMPTZ NOT NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "user_subscription_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  CONSTRAINT "user_subscription_subscription_id_fkey" FOREIGN KEY ("subscription_id") REFERENCES "subscriptions" ("id")
);

CREATE TABLE "subscription_deliveries" (
  "id" bigserial PRIMARY KEY,
  "description" text NULL,
  "user_subscription_id" bigint NOT NULL,
  "delivered_on" TIMESTAMPTZ NOT NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "subscription_deliveries_user_subscription_id_fkey" FOREIGN KEY ("user_subscription_id") REFERENCES "user_subscriptions" ("id")
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "image_url" text[] NOT NULL DEFAULT '{}',
  "product_count" bigint NOT NULL DEFAULT 0,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "description" text NOT NULL,
  "price" decimal(10,2) NOT NULL DEFAULT 0,
  "category_id" bigint NOT NULL,
  "image_url" text[] NOT NULL DEFAULT '{}',
  "stock_quantity" bigint NOT NULL DEFAULT 0,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "products_category_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories" ("id")
);


CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar(255) NOT NULL,
  "user_phone_number" varchar(255) NOT NULL,
  "total_amount" decimal(10,2) NOT NULL,
  "payment_status" boolean NOT NULL,
  "status" varchar(255) NOT NULL,
  "shipping_address" text NULL,
  "deleted_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "quantity" int NOT NULL DEFAULT 1,
  "amount" decimal(10,2) NOT NULL,

  CONSTRAINT "order_item_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "orders" ("id"),
  CONSTRAINT "order_item_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "description" text NULL,
  "order_id" bigint UNIQUE NULL,
  "user_subscription_id" bigint UNIQUE NULL,
  "payment_method" varchar(255) NOT NULL,
  "amount" decimal(10,2) NOT NULL,
  "paid_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),

  CONSTRAINT "payment_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "orders" ("id"),
  CONSTRAINT "payment_user_subscription_id_fkey" FOREIGN KEY ("user_subscription_id") REFERENCES "user_subscriptions" ("id")
);