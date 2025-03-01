package main

type ApiError struct {
	Error string `json:"error"`
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ClerkUser struct {
	Data struct {
		EmailAddresses []struct {
			EmailAddress string `json:"email_address"`
			ID           string `json:"id"`
			LinkedTo     []any  `json:"linked_to"`
			Object       string `json:"object"`
		} `json:"email_addresses"`
		FirstName       string `json:"first_name"`
		ID              string `json:"id"`
		ImageURL        string `json:"image_url"`
		LastName        string `json:"last_name"`
		ProfileImageURL string `json:"profile_image_url"`
		Username        string `json:"username"`
	} `json:"data"`
	Object    string `json:"object"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
}

type NewReleaseBody struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"isPublic"`
}

type Release struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"userId"`
	Imgurl    string `json:"imageUrl"`
	IsPublic  bool   `json:"isPublic"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type NewTrackBody struct {
	Name      string `json:"name"`
	ReleaseId int    `json:"releaseId"`
}

type Track struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ReleaseId int    `json:"releaseId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ReleaseWithTracks struct {
	Release Release `json:"release"`
	Tracks  []Track `json:"tracks"`
}
