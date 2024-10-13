package entities

import (
	"time"
)

type PostDB struct {
	ID       uint32
	AuthorID uint32
	Content  string
	Created  time.Time
	Updated  time.Time
}
