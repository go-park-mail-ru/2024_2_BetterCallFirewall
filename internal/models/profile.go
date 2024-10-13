package models

type FullProfile struct {
	ID        uint32    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Bio       string    `json:"bio"`
	Avatar    Picture   `json:"avatar"`
	Pics      []Picture `json:"pics"`
	Posts     []*Post   `json:"posts"`
}

type ShortProfile struct {
	ID        uint32  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Avatar    Picture `json:"avatar"`
}
