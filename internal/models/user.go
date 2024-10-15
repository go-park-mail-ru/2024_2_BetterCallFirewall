package models

import "net/url"

type Picture url.URL

type User struct {
	ID        uint32  `json:"id"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Avatar    Picture `json:"avatar"`
}