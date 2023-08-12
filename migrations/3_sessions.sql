CREATE TABLE sessions
(
    session_id UUID PRIMARY KEY,
    user_id    UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
);
