package main

type ApiError struct {
	Error string `json:"error"`
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
	IsPublic bool   `json:"is_public"`
	IsSingle bool   `json:"is_single"`
}

type Release struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	Url       string `json:"url"`
	Imgurl    string `json:"image_url"`
	IsPublic  bool   `json:"is_public"`
	IsSingle  bool   `json:"is_single"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Track struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	ReleaseId int    `json:"release_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
