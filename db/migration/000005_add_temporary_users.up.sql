CREATE TABLE "temporary_users" (
  "email" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "token" varchar NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("email", "token")
);

