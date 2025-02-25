-- name: CreateRelease :one
INSERT INTO releases (
  name,
  user_id,
  url,
  imgUrl,
  is_public
) VALUES (
  ?,
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
    url = ?,
    imgUrl = ?,
    is_public = ?,
    updated_at = CURRENT_TIMESTAMP
  WHERE id = ?
RETURNING *;

-- name: GetAllPublicReleases :many
SELECT * FROM releases WHERE is_public = TRUE;

-- name: GetAllPublicReleasesByUser :many
SELECT * FROM releases WHERE is_public = TRUE AND user_id = ?;

-- name: GetAllReleasesByUser :many 
SELECT * FROM releases WHERE user_id = ?;

-- name: DeleteReleaseById :exec
DELETE FROM releases WHERE id = ?;
