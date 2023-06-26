CREATE TABLE "sessions"
(
    "id"                  uuid PRIMARY KEY,
    "username"            varchar     NOT NULL,
    "refresh_token"       varchar     NOT NULL,
    "user_agent"          varchar     NOT NULL,
    "client_ip"           varchar     NOT NULL,
    "is_block"            boolean     NOT NULL DEFAULT false,
    "expires_at"          timestamptz NOT NULL,
    "password_changed_at" timestamptz NOT NULL DEFAULT ('0001-01-01 00:00:00+00')
);

ALTER TABLE "sessions"
    ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

