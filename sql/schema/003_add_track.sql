-- +goose Up
CREATE TABLE tracks (
  id INTEGER UNIQUE PRIMARY KEY,
  name TEXT NOT NULL,
  length INTEGER NOT NULL DEFAULT 0,
  original_file_url TEXT NOT NULL,
  mp3_file_url TEXT NOT NULL,
  release_id INTEGER NOT NULL REFERENCES releases(id) ON DELETE CASCADE,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (release_id) REFERENCES releases(id)
);

-- +goose Down
DROP TABLE tracks;
