package model

// News ...
type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Newses ...
type Newses []*News
