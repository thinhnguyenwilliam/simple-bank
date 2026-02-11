CREATE TABLE sessions (
    id uuid PRIMARY KEY, -- this is id of refresh token
    username varchar NOT NULL,
    refresh_token varchar NOT NULL UNIQUE,
    user_agent varchar NOT NULL,
    client_ip varchar NOT NULL,
    is_blocked boolean NOT NULL DEFAULT false,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_username ON sessions(username);

ALTER TABLE sessions
ADD FOREIGN KEY (username)
REFERENCES users (username);
