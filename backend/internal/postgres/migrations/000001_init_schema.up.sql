CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "email" varchar(50) UNIQUE NOT NULL,
  "address" varchar(255) UNIQUE NOT NULL,
  "phone_number" varchar(50) UNIQUE NOT NULL,
  "refresh_token" text NOT NULL,
  "password" varchar(255) NOT NULL,
  "is_admin" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);