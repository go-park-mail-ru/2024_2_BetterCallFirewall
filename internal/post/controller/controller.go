package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type PostService interface {
	Create(post *models.Post) (uint32, error)
	Get(postID uint32) (*models.Post, error)
	Update(post *models.Post) error
	Delete(postID uint32) error
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)

	ErrorWrongMethod(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
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

func (pc *PostController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	newPost := &models.Post{}

	err := json.NewDecoder(r.Body).Decode(newPost)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
	}

	err = r.ParseMultipartForm(1024 * 1024 * 8 * 5) // 5 M byte
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		pc.responder.ErrorBadRequest(w, myErr.ErrToLargeFile)
	}

	file, _, err := r.FormFile("post")
	defer file.Close()
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
	}

	id, err := pc.service.Create(newPost)
	if err != nil {
		pc.responder.ErrorBadRequest(w, fmt.Errorf("create controller: %w", err))
	}
	newPost.ID = id

	pc.responder.OutputJSON(w, newPost)
}
