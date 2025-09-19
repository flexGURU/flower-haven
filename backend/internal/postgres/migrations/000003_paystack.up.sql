CREATE TABLE "paystack_payments" (
    "id" bigserial PRIMARY KEY,
    "email" varchar(255) NOT NULL,
    "amount" varchar(255) NOT NULL,
    "reference" varchar(255) NOT NULL UNIQUE,
    "status" varchar(50) NOT NULL DEFAULT 'pending',
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "paystack_events" (
    "id" bigserial PRIMARY KEY,
    "event" varchar(255) NOT NULL,
    "data" json NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);