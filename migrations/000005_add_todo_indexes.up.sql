CREATE INDEX IF NOT EXISTS todo_user_id_idx ON todo(user_id);
CREATE INDEX IF NOT EXISTS todo_title_idx ON todo USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS todo_creation_time_idx ON todo (creation_time);
CREATE INDEX IF NOT EXISTS todo_completed_idx ON todo (completed);
