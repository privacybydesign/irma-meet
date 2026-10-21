-- +goose Up
CREATE TABLE irma_meet_sessions (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  irma_session_id text NOT NULL UNIQUE,
  meeting_type text NOT NULL,
  meeting_name text NOT NULL,
  hoster_name text NOT NULL,
  hoster_email_address text NOT NULL,
  invitation_note text NOT NULL DEFAULT '',
  bbb_session_id text NOT NULL,
  bbb_moderator_password text NOT NULL DEFAULT '',
  bbb_attendee_password text NOT NULL DEFAULT '',
  start_time timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE irma_meet_participants (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  irma_session_id text NOT NULL,
  email_address text NOT NULL,
  authentication integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE irma_meet_participants;
DROP TABLE irma_meet_sessions;
