package models

import (
	"time"
)

type Post struct {
	Header    string
	Body      string
	CreatedAt time.Time
}
