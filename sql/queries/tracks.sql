-- name: GetTracksByRelease :many
SELECT * FROM tracks WHERE release_id = ?;

-- name: GetTrackById :one
SELECT * FROM tracks WHERE id = ?;

-- name: CreateTrack :one
INSERT INTO tracks (
  name,
  release_id,
  track_url
) VALUES (
  ?,
  ?,
  ?
) RETURNING *;

-- name: GetTracksByReleaseId :many
SELECT * FROM tracks WHERE release_id = ?;

-- name: GetTracksByUser :many
SELECT * FROM tracks
  INNER JOIN releases ON releases.id = tracks.release_id
  WHERE releases.user_id = ?;

-- name: DeleteTrackById :exec
DELETE FROM tracks WHERE id = ?;
