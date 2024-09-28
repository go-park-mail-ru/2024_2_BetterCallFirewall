package controller

import (
	"errors"
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type PostService interface {
	GetAll() ([]*models.Post, error)
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)

	ErrorNoContent(w http.ResponseWriter, err error)
	ErrorWrongMethod(w http.ResponseWriter, err error)
}

type PostController struct {
	service   PostService
	responder Responder
}

func NewPostController(service PostService, responder Responder) *PostController {
	return &PostController{
		service:   service,
		responder: responder,
	}
}

func (pc *PostController) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	posts, err := pc.service.GetAll()
	if errors.Is(err, myErr.ErrPostEnd) {
		pc.responder.ErrorNoContent(w, err)
		return
	}

	var res []models.Post
	for _, post := range posts {
		res = append(res, *post)
	}
	pc.responder.OutputJSON(w, res)
}
