-- Create "access_token_keys" table
CREATE TABLE "access_token_keys" ("id" uuid NOT NULL, "private_key_encrypted" bytea NOT NULL, "public_pem" character varying NOT NULL, "state" smallint NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "phone_number" character varying NOT NULL, "email" character varying NOT NULL, "username" character varying NOT NULL, "role" smallint NOT NULL, "state" smallint NOT NULL, "first_name" character varying NOT NULL, "last_name" character varying NOT NULL, "password_salt" character varying NOT NULL, "password_hash" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- Create index "users_phone_number_key" to table: "users"
CREATE UNIQUE INDEX "users_phone_number_key" ON "users" ("phone_number");
-- Create index "users_username_key" to table: "users"
CREATE UNIQUE INDEX "users_username_key" ON "users" ("username");
-- Create "refresh_tokens" table
CREATE TABLE "refresh_tokens" ("id" uuid NOT NULL, "is_claimed" boolean NOT NULL DEFAULT false, "token" character varying NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "refresh_tokens_users_refresh_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "refresh_tokens_token_key" to table: "refresh_tokens"
CREATE UNIQUE INDEX "refresh_tokens_token_key" ON "refresh_tokens" ("token");
