package models

//easyjson:json
type Community struct {
	ID               uint32  `json:"id"`
	Name             string  `json:"name"`
	Avatar           Picture `json:"avatar"`
	About            string  `json:"about"`
	CountSubscribers uint32  `json:"count_subscribers"`
	IsAdmin          bool    `json:"is_admin,omitempty"`
	IsFollowed       bool    `json:"is_followed,omitempty"`
}

//easyjson:json
type CommunityCard struct {
	ID         uint32  `json:"id"`
	Name       string  `json:"name"`
	Avatar     Picture `json:"avatar"`
	About      string  `json:"about"`
	IsFollowed bool    `json:"is_followed,omitempty"`
}
