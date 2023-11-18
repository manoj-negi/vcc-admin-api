CREATE TABLE "products" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL,
  "referral_link" varchar,
  "is_active" boolean DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX "product_name" ON "products" ("name");