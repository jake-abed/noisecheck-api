// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: tracks.sql

package database

import (
	"context"
)

const createTrack = `-- name: CreateTrack :one
INSERT INTO tracks (
  name,
  length,
  original_file_url,
  mp3_file_url,
  release_id
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?
) RETURNING id, name, length, original_file_url, mp3_file_url, release_id, created_at, updated_at
`

type CreateTrackParams struct {
	Name            string
	Length          int64
	OriginalFileUrl string
	Mp3FileUrl      string
	ReleaseID       int64
}

func (q *Queries) CreateTrack(ctx context.Context, arg CreateTrackParams) (Track, error) {
	row := q.db.QueryRowContext(ctx, createTrack,
		arg.Name,
		arg.Length,
		arg.OriginalFileUrl,
		arg.Mp3FileUrl,
		arg.ReleaseID,
	)
	var i Track
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Length,
		&i.OriginalFileUrl,
		&i.Mp3FileUrl,
		&i.ReleaseID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteTrackById = `-- name: DeleteTrackById :exec
DELETE FROM tracks WHERE id = ?
`

func (q *Queries) DeleteTrackById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTrackById, id)
	return err
}

const getTrackById = `-- name: GetTrackById :one
SELECT id, name, length, original_file_url, mp3_file_url, release_id, created_at, updated_at FROM tracks WHERE id = ?
`

func (q *Queries) GetTrackById(ctx context.Context, id int64) (Track, error) {
	row := q.db.QueryRowContext(ctx, getTrackById, id)
	var i Track
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Length,
		&i.OriginalFileUrl,
		&i.Mp3FileUrl,
		&i.ReleaseID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTracksByRelease = `-- name: GetTracksByRelease :many
SELECT id, name, length, original_file_url, mp3_file_url, release_id, created_at, updated_at FROM tracks WHERE release_id = ?
`

func (q *Queries) GetTracksByRelease(ctx context.Context, releaseID int64) ([]Track, error) {
	rows, err := q.db.QueryContext(ctx, getTracksByRelease, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Track
	for rows.Next() {
		var i Track
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Length,
			&i.OriginalFileUrl,
			&i.Mp3FileUrl,
			&i.ReleaseID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTracksByReleaseId = `-- name: GetTracksByReleaseId :many
SELECT id, name, length, original_file_url, mp3_file_url, release_id, created_at, updated_at FROM tracks WHERE release_id = ?
`

func (q *Queries) GetTracksByReleaseId(ctx context.Context, releaseID int64) ([]Track, error) {
	rows, err := q.db.QueryContext(ctx, getTracksByReleaseId, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Track
	for rows.Next() {
		var i Track
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Length,
			&i.OriginalFileUrl,
			&i.Mp3FileUrl,
			&i.ReleaseID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTracksByUser = `-- name: GetTracksByUser :many
SELECT tracks.id, tracks.name, length, original_file_url, mp3_file_url, release_id, tracks.created_at, tracks.updated_at, releases.id, releases.name, user_id, imgurl, song_count, is_public, releases.created_at, releases.updated_at FROM tracks
  INNER JOIN releases ON releases.id = tracks.release_id
  WHERE releases.user_id = ?
`

type GetTracksByUserRow struct {
	ID              int64
	Name            string
	Length          int64
	OriginalFileUrl string
	Mp3FileUrl      string
	ReleaseID       int64
	CreatedAt       string
	UpdatedAt       string
	ID_2            int64
	Name_2          string
	UserID          string
	Imgurl          string
	SongCount       int64
	IsPublic        bool
	CreatedAt_2     string
	UpdatedAt_2     string
}

func (q *Queries) GetTracksByUser(ctx context.Context, userID string) ([]GetTracksByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getTracksByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTracksByUserRow
	for rows.Next() {
		var i GetTracksByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Length,
			&i.OriginalFileUrl,
			&i.Mp3FileUrl,
			&i.ReleaseID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Name_2,
			&i.UserID,
			&i.Imgurl,
			&i.SongCount,
			&i.IsPublic,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
