-- name: CreateRelease :one
INSERT INTO releases (
  name,
  user_id,
  url,
  imgUrl,
  is_public,
  is_single
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
)
RETURNING *;
