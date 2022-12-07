ALTER TABLE IF EXISTS user_session
    RENAME COLUMN start_time TO created_at;

ALTER TABLE IF EXISTS user_session
    DROP COLUMN end_time;
