package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"math"
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
	Create(ctx context.Context, post *models.Post) (uint32, error)
	Get(ctx context.Context, postID uint32) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postID uint32) error
	GetBatch(ctx context.Context, lastID uint32) ([]*models.Post, error)
	GetBatchFromFriend(ctx context.Context, userID uint32, lastID uint32) ([]*models.Post, error)
	GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error)
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

	newPost, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	id, err := pc.postService.Create(r.Context(), newPost)
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

	post, err := pc.postService.Get(r.Context(), postID)
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

	post, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		pc.responder.ErrorBadRequest(w, err)
		return
	}

	userID := sess.UserID
	authorID, err := pc.postService.GetPostAuthorID(r.Context(), post.ID)
	if err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err)
			return
		}

		pc.responder.ErrorInternal(w, err)
		return
	}

	if userID != authorID {
		pc.responder.ErrorBadRequest(w, myErr.ErrAccessDenied)
		return
	}

	if err := pc.postService.Update(r.Context(), post); err != nil {
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

	userID := sess.UserID
	authorID, err := pc.postService.GetPostAuthorID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err)
			return
		}

		pc.responder.ErrorInternal(w, err)
		return
	}

	if userID != authorID {
		pc.responder.ErrorBadRequest(w, myErr.ErrAccessDenied)
		return
	}

	if err := pc.postService.Delete(r.Context(), postID); err != nil {
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
		posts  []*models.Post
		err    error
		lastID int
	)

	cookie, err := r.Cookie("postID")
	if err != nil {
		lastID = math.MaxUint32
	} else {
		lastID, err = strconv.Atoi(cookie.Value)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err)
			return
		}
		if lastID < 0 {
			lastID = math.MaxUint32
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

			posts, err = pc.postService.GetBatchFromFriend(r.Context(), sess.UserID, uint32(lastID))
		}
	case "":
		{
			posts, err = pc.postService.GetBatch(r.Context(), uint32(lastID))
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
		Value:   strconv.Itoa(int(posts[len(posts)-1].ID)),
		Path:    "/api/v1/feed/",
		Expires: time.Now().Add(time.Hour),
	}

	http.SetCookie(w, cookie)

	pc.responder.OutputJSON(w, posts)
}

func (pc *PostController) getPostFromBody(r *http.Request) (*models.Post, error) {
	var newPost *models.Post

	if err := json.NewDecoder(r.Body).Decode(newPost); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return nil, err
	}
	newPost.Header.AuthorID = sess.UserID

	defer r.MultipartForm.RemoveAll()
	if err := r.ParseMultipartForm(1024 * 1024 * 8 * 5); err != nil {
		return nil, myErr.ErrToLargeFile
	}

	file, _, err := r.FormFile("file")
	defer file.Close()
	if errors.Is(err, http.ErrMissingFile) {
		return newPost, nil
	}

	if err != nil {
		return nil, err
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
