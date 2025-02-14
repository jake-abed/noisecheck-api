-- +goose Up
CREATE TABLE releases (
  id INTEGER UNIQUE PRIMARY KEY,
  name TEXT NOT NULL,
  song_count INTEGER NOT NULL DEFAULT 0,
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  url TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE releases;
