CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "username" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "ip_addr" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "expired_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY("username") REFERENCES "users" ("username")
