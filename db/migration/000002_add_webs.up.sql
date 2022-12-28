CREATE TABLE "webs" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "url" varchar NOT NULL,
  "title" varchar NOT NULL,
  "thumbnail_url" varchar NOT NULL,
  "html" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "webs" ("user_id");

CREATE UNIQUE INDEX ON "webs" ("user_id", "url");

ALTER TABLE "webs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
