-- +goose Up
CREATE TABLE `users` (
  id string PRIMARY KEY DEFAULT (uuid4()) NOT NULL,
  first_name TEXT,
  last_name TEXT,
  email TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE `users`;
