CREATE TYPE event_enum AS ENUM ('created', 'approved_by_user', 'rejected_by_user', 'signed', 'sent');

CREATE TABLE IF NOT EXISTS task_events (
  task_uuid uuid,
  event     event_enum,
  user_uuid uuid,
  timestamp timestamptz
);