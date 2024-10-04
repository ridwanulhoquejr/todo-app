CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  name text NOT NULL,
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  activated bool NOT NULL DEFAULT (false),
  creation_time timestamp(0) with time zone NOT NULL DEFAULT (now())
);
