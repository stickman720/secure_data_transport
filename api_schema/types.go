package api_schema




type TeamInput struct {
	Title string `json:"title"`
	UserCount int `json:"user-count"`
	RateLimitPerUser int `json:"rate-limit-per-user"`
	StorageSpace int `json:"storage-space"`
	Transfer int `json:"transfear"`
}
