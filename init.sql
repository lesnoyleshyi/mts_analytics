CREATE TYPE event_enum AS ENUM ('created', 'sent_to', 'approved_by', 'rejected_by', 'signed', 'sent');

CREATE TABLE IF NOT EXISTS task_events (
  task_uuid uuid,
  event     event_enum,
  user_uuid uuid DEFAULT '00000000-0000-0000-0000-000000000000',
  timestamp timestamptz,
  CONSTRAINT task_events_PK PRIMARY KEY (task_uuid, event, user_uuid, timestamp)
);

-- INSERT INTO task_events (task_uuid, event, user_uuid, timestamp) VALUES
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'created', '0a75fc43-584a-46a6-8e3c-686edbea43f8', '2022-06-10 10:38:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'approved_by', '97364b12-0dcf-487f-8a03-9051d5c1d8ba', '2022-06-10 11:00:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'approved_by', 'c158314b-a043-4e25-8bd0-4de99bfd4601', '2022-06-10 11:05:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'approved_by', '469481c6-0376-4a9f-a2e0-8d72b99bde67', '2022-06-10 12:30:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'signed', NULL, '2022-06-10 12:30:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'sent', NULL, '2022-06-10 12:31:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'created', 'eafe0883-5aa0-4f8a-850f-ed002b56c131', '2022-06-10 10:38:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'approved_by', '0a75fc43-584a-46a6-8e3c-686edbea43f8', '2022-06-11 11:00:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'rejected_by', 'c158314b-a043-4e25-8bd0-4de99bfd4601', '2022-06-11 11:05:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'approved_by', '97364b12-0dcf-487f-8a03-9051d5c1d8ba', '2022-06-10 12:30:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'created', 'c158314b-a043-4e25-8bd0-4de99bfd4601', '2022-05-09 11:38:41+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'approved_by', '0a75fc43-584a-46a6-8e3c-686edbea43f8', '2022-05-12 11:00:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'approved_by', '97364b12-0dcf-487f-8a03-9051d5c1d8ba', '2022-05-20 11:05:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'approved_by', '469481c6-0376-4a9f-a2e0-8d72b99bde67', '2022-05-26 12:30:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'approved_by', 'eafe0883-5aa0-4f8a-850f-ed002b56c131', '2022-06-07 12:30:40+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'signed', NULL, '2022-06-07 12:31:00+00'),
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'sent', NULL, '2022-06-07 12:31:10+00');



