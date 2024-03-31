CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now ()),
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "is_email_verified" bool NOT NULL DEFAULT 'false'
);

CREATE TABLE "conversation" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now ())
);

CREATE TABLE "message" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "message" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now ()),
  "conversation_id" bigint NOT NULL
);

CREATE TABLE "conversation_users" (
  "id" bigserial PRIMARY KEY,
  "conversation_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now ())
);

CREATE INDEX ON "message" ("user_id");

CREATE INDEX ON "message" ("conversation_id");

CREATE INDEX ON "message" ("user_id", "conversation_id");

CREATE INDEX ON "conversation_users" ("user_id");

CREATE INDEX ON "conversation_users" ("conversation_id");

CREATE INDEX ON "conversation_users" ("user_id", "conversation_id");

ALTER TABLE "message" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "message" ADD FOREIGN KEY ("conversation_id") REFERENCES "conversation" ("id");

ALTER TABLE "conversation_users" ADD FOREIGN KEY ("conversation_id") REFERENCES "conversation" ("id") ON DELETE CASCADE;

ALTER TABLE "conversation_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "conversation_users" ADD CONSTRAINT "conversation_user_key" UNIQUE ("conversation_id", "user_id")
