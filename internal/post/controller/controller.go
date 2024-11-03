package controller

import (
	"context"
	"errors"
	"fmt"
	"image"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"

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
	LogError(err error, requestId string)
}

type FileService interface {
	Download(ctx context.Context, file multipart.File, postID, profileID uint32) error
	GetPostPicture(ctx context.Context, postID uint32) *models.Picture
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
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

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

	err = pc.saveFileFromMultipart(r, id)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	newPost.ID = id

	pc.responder.OutputJSON(w, newPost, reqID)
}

func (pc *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

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

	post.PostContent.File = pc.fileService.GetPostPicture(r.Context(), postID)

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) Update(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

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

	err = pc.saveFileFromMultipart(r, post.ID)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok   = r.Context().Value("requestID").(string)
		postID, err = getIDFromQuery(r)
	)

	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

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
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		section   = r.URL.Query().Get("section")
		lastID    = r.URL.Query().Get("id")
		posts     []*models.Post
		intLastID uint64
		err       error
	)
	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if lastID == "" {
		intLastID = math.MaxInt32
	} else {
		intLastID, err = strconv.ParseUint(lastID, 10, 32)
		if err != nil {
			pc.responder.ErrorBadRequest(w, myErr.ErrInvalidQuery, reqID)
			return
		}
	}

	switch section {
	case "friend":
		{
			sess, errSession := models.SessionFromContext(r.Context())
			if errSession != nil {
				pc.responder.ErrorBadRequest(w, errSession, reqID)
				return
			}

			posts, err = pc.postService.GetBatchFromFriend(r.Context(), sess.UserID, uint32(intLastID))
		}
	case "":
		{
			posts, err = pc.postService.GetBatch(r.Context(), uint32(intLastID))
		}
	default:
		pc.responder.ErrorBadRequest(w, myErr.ErrInvalidQuery, reqID)
		return
	}

	if err != nil && !errors.Is(err, myErr.ErrNoMoreContent) && !errors.Is(err, myErr.ErrAnotherService) {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if errors.Is(err, myErr.ErrAnotherService) {
		pc.responder.LogError(myErr.ErrAnotherService, reqID)
	}

	for _, p := range posts {
		p.PostContent.File = pc.fileService.GetPostPicture(r.Context(), p.ID)
	}

	if errors.Is(err, myErr.ErrNoMoreContent) {
		pc.responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	pc.responder.OutputJSON(w, posts, reqID)
}

func (pc *PostController) getPostFromBody(r *http.Request) (*models.Post, error) {
	var newPost *models.Post

	err := r.ParseMultipartForm(10 << 20) // 10Mbyte
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		return nil, myErr.ErrToLargeFile
	}

	text := r.Form.Get("text")
	newPost.PostContent.Text = text

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return nil, err
	}
	newPost.Header.AuthorID = sess.UserID

	return newPost, nil
}

func (pc *PostController) saveFileFromMultipart(r *http.Request, postID uint32) error {
	if r.MultipartForm == nil {
		return nil
	}

	err := r.ParseMultipartForm(10 << 20) // 10Mbyte
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		return myErr.ErrToLargeFile
	}

	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		return myErr.ErrWrongMultipartForm
	}

	_, format, err := image.Decode(file)
	if err != nil {
		return myErr.ErrWrongMultipartForm
	}

	if _, ok := fileFormat[format]; !ok {
		return myErr.ErrWrongFiletype
	}

	err = pc.fileService.Download(r.Context(), file, postID, 0)

	return err
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
