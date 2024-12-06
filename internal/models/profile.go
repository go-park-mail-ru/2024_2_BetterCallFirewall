package models

type FullProfile struct {
	ID             uint32    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Bio            string    `json:"bio"`
	IsAuthor       bool      `json:"is_author"`
	IsFriend       bool      `json:"is_friend"`
	IsSubscriber   bool      `json:"is_subscriber"`
	IsSubscription bool      `json:"is_subscription"`
	Avatar         Picture   `json:"avatar"`
	Pics           []Picture `json:"pics"`
	Posts          []*Post   `json:"posts"`
}

type ShortProfile struct {
	ID             uint32  `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	IsAuthor       bool    `json:"is_author"`
	IsFriend       bool    `json:"is_friend"`
	IsSubscriber   bool    `json:"is_subscriber"`
	IsSubscription bool    `json:"is_subscription"`
	Avatar         Picture `json:"avatar"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
