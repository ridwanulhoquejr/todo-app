CREATE TABLE IF NOT EXISTS todo (
  id bigserial PRIMARY KEY,
  user_id integer NOT NULL,
  title text NOT NULL,
  description text,
  completed bool NOT NULL DEFAULT (false),
  version integer NOT NULL DEFAULT 0,
  creation_time timestamp(0) with time zone NOT NULL DEFAULT (now())
);

CREATE INDEX ON todo (title);