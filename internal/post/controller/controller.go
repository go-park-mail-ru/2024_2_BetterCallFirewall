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
	CheckUserAccess(userID uint32, postID uint32) (bool, error)
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)

	ErrorWrongMethod(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
}

type FileService interface {
	Save(file multipart.File, fileHeader *multipart.FileHeader) (*models.Picture, error)
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

	newPost, err := pc.getPost(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

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

	post, err := pc.getPost(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	authorID := sess.UserID
	ok, err := pc.service.CheckUserAccess(authorID, post.ID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	if !ok {
		pc.responder.ErrorBadRequest(w, errors.New("access denied"))
	}

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
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	authorID := sess.UserID
	ok, err := pc.service.CheckUserAccess(authorID, postID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	if !ok {
		pc.responder.ErrorBadRequest(w, errors.New("access denied"))
	}

	if err := pc.service.Delete(postID); err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	pc.responder.OutputJSON(w, postID)
}

func (pc *PostController) getPost(r *http.Request) (*models.Post, error) {
	var newPost *models.Post

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return nil, err
	}
	newPost.AuthorID = sess.UserID

	defer r.MultipartForm.RemoveAll()
	if err := r.ParseMultipartForm(1024 * 1024 * 8 * 5); err != nil {
		return nil, myErr.ErrToLargeFile
	}

	file, fileHeader, err := r.FormFile("file")
	defer file.Close()
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return nil, err
	}

	if errors.Is(err, http.ErrMissingFile) {
		return newPost, nil
	}

	pic, err := pc.fileService.Save(file, fileHeader)
	if err != nil {
		return nil, err
	}
	newPost.PostContent.File = pic

	return newPost, nil
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
