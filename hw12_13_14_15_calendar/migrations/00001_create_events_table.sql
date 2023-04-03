-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL SECOND(0) NOT NULL,
    description TEXT,
    user_id UUID NOT NULL,
    remind_date INTERVAL SECOND(0)
);

CREATE INDEX event_date_idx
    ON events (date);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS events;