-- name: GetTracksByRelease :many
SELECT * FROM tracks WHERE release_id = ?;

-- name: GetTrackById :one
SELECT * FROM tracks WHERE id = ?;

-- name: AddTrack :one
INSERT INTO tracks (
  name,
  url,
  release_id
) VALUES (
  ?,
  ?,
  ?
) RETURNING *;

-- name: GetTracksByUser :many
SELECT * FROM tracks
  INNER JOIN releases ON releases.id = tracks.release_id
  WHERE releases.user_id = ?;

-- name: DeleteTrackById :exec
DELETE FROM tracks WHERE id = ?;
