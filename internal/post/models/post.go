package models

type Post struct {
	Header    string `json:"header"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}
