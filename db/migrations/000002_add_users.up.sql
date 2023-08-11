-- TODO: username shouldn't be a primary key, use email instead.
CREATE TABLE "users" (
    "email" varchar UNIQUE NOT NULL,
    "username" varchar PRIMARY KEY,
    "password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY("owner_name") REFERENCES "users"("username");

-- user can't have duplicate accounts(user account) with duplicate currencies
-- CREATE UNIQUE INDEX ON "accounts" ("owner_name", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner_name", "currency")
