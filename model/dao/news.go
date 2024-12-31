package dao

// News ...
type News struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"articles"`
}
