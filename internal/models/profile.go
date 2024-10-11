package models

type Picture Picture

type Profile struct {
	ID        uint32    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Bio       string    `json:"bio"`
	Avatar    Picture   `json:"avatar"`
	Pics      []Picture `json:"pics"`
	Posts     []string  //TODO make posts instead
}
