package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

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
	OutputJSON(w http.ResponseWriter, data any, requestId string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestId string)

	ErrorInternal(w http.ResponseWriter, err error, requestId string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestId string)
}

type FileService interface {
	Upload(file multipart.File) (*models.Picture, error)
	GetPostPicture(postID uint32) *models.Picture
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
	reqID := r.Context().Value("requestID").(string)

	newPost, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	id, err := pc.postService.Create(r.Context(), newPost)
	if err != nil {
		pc.responder.ErrorInternal(w, fmt.Errorf("create controller: %w", err), reqID)
		return
	}
	newPost.ID = id

	pc.responder.OutputJSON(w, newPost, reqID)
}

func (pc *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)

	postID, err := getIDFromQuery(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	post, err := pc.postService.Get(r.Context(), postID)
	if err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if !errors.Is(err, myErr.ErrAnotherService) {
			pc.responder.ErrorInternal(w, err, reqID)
			return
		}
	}

	post.PostContent.File = pc.fileService.GetPostPicture(postID)

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) Update(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)

	post, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	userID := post.Header.AuthorID
	authorID, err := pc.postService.GetPostAuthorID(r.Context(), post.ID)
	if err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}

		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if userID != authorID {
		pc.responder.ErrorBadRequest(w, myErr.ErrAccessDenied, reqID)
		return
	}

	if err := pc.postService.Update(r.Context(), post); err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	postID, err := getIDFromQuery(r)

	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	userID := sess.UserID
	authorID, err := pc.postService.GetPostAuthorID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}

		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if userID != authorID {
		pc.responder.ErrorBadRequest(w, myErr.ErrAccessDenied, reqID)
		return
	}

	if err := pc.postService.Delete(r.Context(), postID); err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, postID, reqID)
}

func (pc *PostController) GetBatchPosts(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	section := r.URL.Query().Get("section")

	var (
		posts  []*models.Post
		err    error
		lastID int
	)

	cookie, err := r.Cookie("postID")
	if err != nil {
		lastID = math.MaxInt32
	} else {
		lastID, err = strconv.Atoi(cookie.Value)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if lastID == 0 {
			pc.responder.OutputNoMoreContentJSON(w, reqID)
			return
		}
		if lastID < 0 {
			lastID = math.MaxInt32
		}
	}

	switch section {
	case "friend":
		{
			sess, errSession := models.SessionFromContext(r.Context())
			if errSession != nil {
				pc.responder.ErrorBadRequest(w, err, reqID)
				return
			}

			posts, err = pc.postService.GetBatchFromFriend(r.Context(), sess.UserID, uint32(lastID))
		}
	case "":
		{
			posts, err = pc.postService.GetBatch(r.Context(), uint32(lastID))
		}
	default:
		pc.responder.ErrorBadRequest(w, errors.New("invalid query params"), reqID)
		return
	}

	if err != nil && !errors.Is(err, myErr.ErrNoMoreContent) && !errors.Is(err, myErr.ErrAnotherService) {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if errors.Is(err, myErr.ErrAnotherService) {
		log.Printf("req: %s: erranother: %v", reqID, err)
	}

	for _, p := range posts {
		p.PostContent.File = pc.fileService.GetPostPicture(p.ID)
	}

	var newLastID string

	if errors.Is(err, myErr.ErrNoMoreContent) {
		newLastID = "0"
	} else {
		newLastID = strconv.Itoa(int(posts[len(posts)-1].ID))
	}

	cookie = &http.Cookie{
		Name:    "postID",
		Path:    "/api/v1/feed",
		Value:   newLastID,
		Expires: time.Now().Add(time.Hour),
	}

	http.SetCookie(w, cookie)

	pc.responder.OutputJSON(w, posts, reqID)
}

func (pc *PostController) getPostFromBody(r *http.Request) (*models.Post, error) {
	var newPost *models.Post

	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return nil, err
	}
	newPost.Header.AuthorID = sess.UserID

	return newPost, nil
}

func getIDFromQuery(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)

	id := vars["id"]
	if id == "" {
		return 0, errors.New("id is empty")
	}

	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(uid), nil
}
