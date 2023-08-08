CREATE TABLE tasks (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  status BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL,
  edited_at TIMESTAMP NOT NULL
);


