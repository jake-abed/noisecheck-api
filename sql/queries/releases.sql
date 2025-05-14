-- name: CreateRelease :one
INSERT INTO releases (
  name,
  user_id,
  imgUrl,
  is_public
) VALUES (
  ?,
  ?,
  ?,
  ?
)
RETURNING *;

-- name: GetReleaseById :one
SELECT * FROM releases WHERE id = ?;

-- name: UpdateRelease :one
UPDATE releases
  SET name = ?,
    imgUrl = ?,
    is_public = ?,
    updated_at = CURRENT_TIMESTAMP
  WHERE id = ?
RETURNING *;

-- name: GetPublicReleases :many
SELECT releases.*, users.username FROM releases
  INNER JOIN users ON users.id = releases.user_id
  WHERE is_public = TRUE
  ORDER BY releases.created_at DESC
  LIMIT 20 OFFSET ?;

-- name: GetAllPublicReleases :many
SELECT * FROM releases WHERE is_public = TRUE;

-- name: GetAllPublicReleasesByUser :many
SELECT releases.*, users.username FROM releases
  INNER JOIN users ON users.id = releases.user_id
  WHERE is_public = TRUE AND user_id = ?
  ORDER BY releases.created_at DESC
  LIMIT 20 OFFSET ?;

-- name: GetAllReleasesByUser :many 
SELECT releases.*, users.username FROM releases
  INNER JOIN users ON users.id = releases.user_id
  WHERE user_id = ?
  ORDER BY releases.created_at DESC
  LIMIT 20 OFFSET ?;

-- name: DeleteReleaseById :exec
DELETE FROM releases WHERE id = ?;
