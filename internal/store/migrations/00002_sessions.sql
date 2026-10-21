-- +goose Up
-- Server-side HTTP sessions (alexedwards/scs postgresstore).
CREATE TABLE sessions (
  token text PRIMARY KEY,
  data bytea NOT NULL,
  expiry timestamptz NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose Down
DROP TABLE sessions;
