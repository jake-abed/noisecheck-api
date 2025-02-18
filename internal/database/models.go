// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

type Release struct {
	ID        int64
	Name      string
	UserID    string
	Url       string
	Imgurl    string
	SongCount int64
	IsPublic  bool
	IsSingle  bool
	CreatedAt string
	UpdatedAt string
}

type Track struct {
	ID        int64
	Name      string
	Url       string
	ReleaseID int64
	CreatedAt string
	UpdatedAt string
}

type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt string
	UpdatedAt string
}
