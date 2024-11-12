package models

type Community struct {
	ID               uint32  `json:"id"`
	Name             string  `json:"name"`
	Avatar           Picture `json:"avatar"`
	About            string  `json:"about"`
	CountSubscribers uint32  `json:"count_subscribers"`
}

type CommunityCard struct {
	ID     uint32  `json:"id"`
	Name   string  `json:"name"`
	Avatar Picture `json:"avatar"`
	About  string  `json:"about"`
}
