-- migrations/000007_product_seeder.up.sql

INSERT INTO "products" ("name", "referral_link")
VALUES
  ('roottester', 'https://www.roottester.com'),
  ('outline', 'https://www.outline.com'),
  ('boringant','https://www.boringant.com');
