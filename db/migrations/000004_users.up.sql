CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" varchar NOT NULL,
  "role_id" int,
  "api_key" varchar NOT NULL,
  "client_id" varchar NOT NULL,
  "country_code" int,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "validation_token" varchar,
  "mobile" varchar,
  "referral_code" varchar,
  "product_id" int NOT NULL,
  "total_invitees" int DEFAULT 0,
  "successful_referral" int DEFAULT 0,
  "is_active" int DEFAULT 1,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "user_roles" ("id");
ALTER TABLE "users" ADD FOREIGN KEY ("country_code") REFERENCES "countries" ("code");
ALTER TABLE "users" ADD CONSTRAINT "users_email_product_id_key" UNIQUE ("email", "product_id");
CREATE INDEX "user_email_index" ON "users" ("email", "username");