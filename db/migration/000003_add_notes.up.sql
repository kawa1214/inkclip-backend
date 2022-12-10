CREATE TABLE "notes" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "title" varchar NOT NULL,
  "content" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "note_webs" (
  "note_id" uuid NOT NULL,
  "web_id" uuid NOT NULL,
  PRIMARY KEY ("note_id", "web_id")
);


CREATE UNIQUE INDEX ON "users" ("email");

CREATE INDEX ON "webs" ("user_id");

CREATE UNIQUE INDEX ON "webs" ("user_id", "url");

CREATE INDEX ON "notes" ("user_id");

CREATE INDEX ON "note_webs" ("note_id");

ALTER TABLE "webs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "notes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "note_webs" ADD FOREIGN KEY ("note_id") REFERENCES "notes" ("id") ON DELETE CASCADE;

ALTER TABLE "note_webs" ADD FOREIGN KEY ("web_id") REFERENCES "webs" ("id") ON DELETE CASCADE;
