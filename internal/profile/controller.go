package profile

import (
	"net/http"
)

type ProfileHandler interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
	GetProfileById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)

	SendFriendReq(w http.ResponseWriter, r *http.Request)
	AcceptFriendReq(w http.ResponseWriter, r *http.Request)
	RemoveFromFriends(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	GetAllFriends(w http.ResponseWriter, r *http.Request)
	GetAllSubs(w http.ResponseWriter, r *http.Request)
	GetAllSubscriptions(w http.ResponseWriter, r *http.Request)
}
