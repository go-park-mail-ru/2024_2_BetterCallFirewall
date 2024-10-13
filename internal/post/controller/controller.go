package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

var fileFormat = map[string]struct{}{
	"jpeg": {},
	"jpg":  {},
	"png":  {},
}

type PostService interface {
	Create(post *models.Post) (uint32, error)
	Get(postID uint32) (*models.Post, error)
	Update(post *models.Post) error
	Delete(postID uint32) error
	GetBatch(lastID uint32, newRequest bool) ([]*models.Post, error)
	GetBatchFromFriend(userID uint32, lastID uint32, newRequest bool) ([]*models.Post, error)
	CheckUserAccess(userID uint32, postID uint32) (bool, error)
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)
	OutputNoMoreContentJSON(w http.ResponseWriter, data any)

	ErrorInternal(w http.ResponseWriter, err error)
	ErrorWrongMethod(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
}

type FileService interface {
	Upload(file multipart.File) (*models.Picture, error)
	GetPostPicture(postID uint32) (*models.Picture, error)
}

type PostController struct {
	postService PostService
	responder   Responder
	fileService FileService
}

func NewPostController(service PostService, responder Responder, fileService FileService) *PostController {
	return &PostController{
		postService: service,
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

	id, err := pc.postService.Create(newPost)
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

	post, err := pc.postService.Get(postID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	post.PostContent.File, err = pc.fileService.GetPostPicture(postID)
	if err != nil {
		pc.responder.ErrorInternal(w, err)
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
	ok, err := pc.postService.CheckUserAccess(authorID, post.ID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}
	if !ok {
		pc.responder.ErrorBadRequest(w, errors.New("access denied"))
	}

	if err := pc.postService.Update(post); err != nil {
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
	ok, err := pc.postService.CheckUserAccess(authorID, postID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	if !ok {
		pc.responder.ErrorBadRequest(w, errors.New("access denied"))
	}

	if err := pc.postService.Delete(postID); err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	pc.responder.OutputJSON(w, postID)
}

func (pc *PostController) GetBatchPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pc.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	section := r.URL.Query().Get("section")

	var (
		posts      []*models.Post
		err        error
		lastID     int
		newRequest = false
	)

	cookie, err := r.Cookie("postID")
	if err != nil {
		lastID = 0
		newRequest = true
	} else {
		lastID, err = strconv.Atoi(cookie.Value)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err)
			return
		}
		if lastID < 0 {
			lastID = 0
			newRequest = true
		}
	}

	switch section {
	case "friend":
		{
			sess, errSession := models.SessionFromContext(r.Context())
			if errSession != nil {
				pc.responder.ErrorBadRequest(w, err)
				return
			}

			posts, err = pc.postService.GetBatchFromFriend(sess.UserID, uint32(lastID), newRequest)
		}
	case "":
		{
			posts, err = pc.postService.GetBatch(uint32(lastID), newRequest)
		}
	default:
		pc.responder.ErrorBadRequest(w, errors.New("invalid query params"))
		return
	}

	for _, p := range posts {
		p.PostContent.File, err = pc.fileService.GetPostPicture(p.ID)
		if err != nil {
			pc.responder.ErrorInternal(w, err)
			return
		}
	}

	if errors.Is(err, myErr.ErrNoMoreContent) {
		pc.responder.OutputNoMoreContentJSON(w, posts)
		return
	}

	if err != nil {
		pc.responder.ErrorInternal(w, err)
		return
	}

	cookie = &http.Cookie{
		Name:    "postID",
		Value:   strconv.Itoa(int(lastID)),
		Path:    "/api/v1/feed/",
		Expires: time.Now().Add(time.Hour),
	}

	http.SetCookie(w, cookie)

	pc.responder.OutputJSON(w, posts)
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

	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return nil, err
	}

	if errors.Is(err, http.ErrMissingFile) {
		return newPost, nil
	}

	_, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if _, ok := fileFormat[format]; !ok {
		return nil, myErr.ErrWrongFiletype
	}

	pic, err := pc.fileService.Upload(file)
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
