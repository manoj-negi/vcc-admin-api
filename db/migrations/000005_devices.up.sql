CREATE TABLE "devices" (
  "id" SERIAL PRIMARY KEY,
  "os" varchar,
  "os_version" varchar,
  "manufacturer" varchar,
  "model" varchar,
  "device_token" varchar UNIQUE,
  "user_id" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT NOW(),
  "updated_at" timestamp NOT NULL DEFAULT NOW()
);

ALTER TABLE "devices" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");