-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id VARCHAR ( 50 ) PRIMARY KEY,
    title TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    duration TIME NOT NULL,
    description TEXT,
    user_id VARCHAR ( 50 ) NOT NULL,
    remind_date TIME
);

CREATE INDEX event_date_idx
    ON events (date);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS events;