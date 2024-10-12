package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
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

type FileService interface {
	Save(file multipart.File) (models.Picture, error)
}

type PostController struct {
	service     PostService
	responder   Responder
	fileService FileService
}

func NewPostController(service PostService, responder Responder, fileService FileService) *PostController {
	return &PostController{
		service:     service,
		responder:   responder,
		fileService: fileService,
	}
}

func (pc *PostController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	var newPost *models.Post

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	defer r.Body.Close()

	defer r.MultipartForm.RemoveAll()
	if err := r.ParseMultipartForm(1024 * 1024 * 8 * 5); err != nil {
		pc.responder.ErrorBadRequest(w, myErr.ErrToLargeFile)
		return
	}

	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	pic, err := pc.fileService.Save(file)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	newPost.PostContent.File = pic
	id, err := pc.service.Create(newPost)
	if err != nil {
		pc.responder.ErrorBadRequest(w, fmt.Errorf("create controller: %w", err))
		return
	}
	newPost.ID = id

	pc.responder.OutputJSON(w, newPost)
}

func (pc *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	postID, err := getIDFromQuery(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	post, err := pc.service.Get(postID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	pc.responder.OutputJSON(w, post)
}

func (pc *PostController) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	var post *models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	defer r.Body.Close()

	if err := pc.service.Update(post); err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	pc.responder.OutputJSON(w, post)
}

func (pc *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	postID, err := getIDFromQuery(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	if err := pc.service.Delete(postID); err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	pc.responder.OutputJSON(w, postID)
}

func getIDFromQuery(r *http.Request) (uint32, error) {
	id := r.URL.Query().Get("id")

	if id == "" {
		return 0, errors.New("wrong id")
	}

	postID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("wrong id")
	}

	return uint32(postID), nil
}
