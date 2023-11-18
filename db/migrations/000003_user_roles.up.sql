CREATE TABLE "user_roles" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar NOT NULL
);
CREATE INDEX "role_name" ON "user_roles" ("name");