ALTER TABLE "products" DROP COLUMN "has_stems";
ALTER TABLE "products" DROP COLUMN "is_message_card";
ALTER TABLE "products" DROP COLUMN "is_flowers";
ALTER TABLE "products" DROP COLUMN "is_add_on";

DROP TABLE IF EXISTS product_stems;