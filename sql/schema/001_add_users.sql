-- +goose Up
CREATE TABLE users (
  id string PRIMARY KEY DEFAULT (uuid4()) NOT NULL,
  username TEXT,
  email TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE `users`;
