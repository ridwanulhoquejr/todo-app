CREATE TABLE IF NOT EXISTS token(
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);