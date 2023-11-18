CREATE TABLE "verify_emails" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL,
  "email" varchar,
  "secret_code" varchar,
  "is_used" boolean,
  "expires_at" timestamp NOT NULL DEFAULT (now()),
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
CREATE INDEX idx_verify_emails_user_id ON verify_emails (user_id);