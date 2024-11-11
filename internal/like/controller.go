package like

import (
	"net/http"
)

type ReactionController interface {
	SetLikeToPost(w http.ResponseWriter, r *http.Request)
	SetLikeToComment(w http.ResponseWriter, r *http.Request)
	SetLikeToFile(w http.ResponseWriter, r *http.Request)
	DeleteLikeFromPost(w http.ResponseWriter, r *http.Request)
	DeleteLikeFromComment(w http.ResponseWriter, r *http.Request)
	DeleteLikeFromFile(w http.ResponseWriter, r *http.Request)
}
