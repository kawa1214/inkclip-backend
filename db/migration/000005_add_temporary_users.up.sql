CREATE TABLE "temporary_users" (
  "email" varchar PRIMARY KEY NOT NULL,
  "hashed_password" varchar NOT NULL,
  "token" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "temporary_users" ("email");