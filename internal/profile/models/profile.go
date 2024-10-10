package models

import (
	"net/url"
)

type Picture url.URL

type Profile struct {
	ID        uint32
	FirstName string
	LastName  string
	Bio       string
	Avatar    Picture
	Pics      []Picture
}
