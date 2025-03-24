CREATE TABLE events
(
    id            SERIAL PRIMARY KEY,
    title         TEXT      NOT NULL,
    start_time    TIMESTAMP NOT NULL,
    duration      INTERVAL  NOT NULL,
    description   TEXT,
    user_id       TEXT      NOT NULL,
    notify_before INTERVAL
);