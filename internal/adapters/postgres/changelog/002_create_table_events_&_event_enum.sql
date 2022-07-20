-- +goose Up
CREATE TYPE app.event_enum AS ENUM ('created', 'sent_to', 'approved_by', 'rejected_by', 'signed', 'sent');

CREATE TABLE IF NOT EXISTS app.task_events (
    task_uuid uuid,
    event     app.event_enum,
    user_uuid uuid DEFAULT '00000000-0000-0000-0000-000000000000',
    timestamp timestamptz,
    CONSTRAINT task_events_PK PRIMARY KEY (task_uuid, event, user_uuid, timestamp)
);

-- +goose Down
DROP TABLE app.task_events;

DROP TYPE app.event_enum;
