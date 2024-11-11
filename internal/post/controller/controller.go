package controller

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

var fileFormat = map[string]struct{}{
	"image/jpeg": {},
	"image/jpg":  {},
	"image/png":  {},
	"image/webp": {},
}

type PostService interface {
	Create(ctx context.Context, post *models.Post) (uint32, error)
	Get(ctx context.Context, postID uint32) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postID uint32) error
	GetBatch(ctx context.Context, lastID uint32) ([]*models.Post, error)
	GetBatchFromFriend(ctx context.Context, userID uint32, lastID uint32) ([]*models.Post, error)
	GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error)

	GetCommunityPost(ctx context.Context, communityID, lastID uint32) ([]*models.Post, error)
	CreateCommunityPost(ctx context.Context, post *models.Post) (uint32, error)
	CheckAccessToCommunity(ctx context.Context, userID uint32, postID uint32) bool
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
	UpdatePostFile(ctx context.Context, file multipart.File, postID uint32) error
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

	newPost, file, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	id, err := pc.postService.Create(r.Context(), newPost)
	if err != nil {
		pc.responder.ErrorInternal(w, fmt.Errorf("create controller: %w", err), reqID)
		return
	}

	if file != nil {
		defer file.Close()
		err = pc.fileService.Download(r.Context(), file, id, 0)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
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
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = getIDFromQuery(r)
	)

	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	post, file, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	post.ID = id
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

	if file != nil {
		defer file.Close()
		err = pc.fileService.UpdatePostFile(r.Context(), file, post.ID)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
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
		posts     []*models.Post
		intLastID uint64
		err       error
	)

	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	intLastID, err = getLastID(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
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

func (pc *PostController) GetCommunityPosts(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	lastID, err := getLastID(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	id := r.URL.Query().Get("community_id")
	if id == "" {
		pc.responder.ErrorBadRequest(w, myErr.ErrInvalidQuery, reqID)
		return
	}
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		pc.responder.ErrorBadRequest(w, myErr.ErrInvalidQuery, reqID)
	}

	posts, err := pc.postService.GetCommunityPost(r.Context(), uint32(uid), uint32(lastID))
	if err != nil {
		if errors.Is(err, myErr.ErrNoMoreContent) {
			pc.responder.OutputNoMoreContentJSON(w, reqID)
			return
		}
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, posts, reqID)
}

func (pc *PostController) CreatePostInCommunity(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	post, file, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !checkAccess(post.Header.AuthorID, post.Header.CommunityID) {
		pc.responder.ErrorBadRequest(w, myErr.ErrAccessDenied, reqID)
		return
	}

	id, err := pc.postService.CreateCommunityPost(r.Context(), post)
	if err != nil {
		pc.responder.ErrorInternal(w, fmt.Errorf("create controller: %w", err), reqID)
		return
	}

	if file != nil {
		defer file.Close()
		err = pc.fileService.Download(r.Context(), file, id, 0)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
	}
	post.ID = id

	pc.responder.OutputJSON(w, post, reqID)

}

func (pc *PostController) UpdatePostInCommunity(w http.ResponseWriter, r *http.Request) {
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

	if !pc.postService.CheckAccessToCommunity(r.Context(), sess.UserID, postID) {
		pc.responder.ErrorBadRequest(w, myErr.ErrAccessDenied, reqID)
		return
	}

	post, file, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	post.ID = postID

	if err := pc.postService.Update(r.Context(), post); err != nil {
		if errors.Is(err, myErr.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if file != nil {
		defer file.Close()
		err = pc.fileService.UpdatePostFile(r.Context(), file, post.ID)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
	}

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) DeletePostInCommunity(w http.ResponseWriter, r *http.Request) {
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

	if !pc.postService.CheckAccessToCommunity(r.Context(), sess.UserID, postID) {
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

func (pc *PostController) getPostFromBody(r *http.Request) (*models.Post, multipart.File, error) {
	var newPost models.Post

	if r.MultipartForm == nil {
		return nil, nil, errors.New("multipart form required")
	}

	err := r.ParseMultipartForm(10 << 20) // 10Mbyte
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		return nil, nil, myErr.ErrToLargeFile
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		file = nil
	} else {
		format := header.Header.Get("Content-Type")
		if _, ok := fileFormat[format]; !ok {
			return nil, nil, myErr.ErrWrongFiletype
		}
	}

	text := r.Form.Get("text")
	newPost.PostContent.Text = text
	communityID := r.Form.Get("community_id")
	intCommunityID, err := strconv.ParseUint(communityID, 10, 32)
	if err != nil {
		intCommunityID = 0
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return nil, nil, err
	}
	newPost.Header.AuthorID = sess.UserID
	newPost.Header.CommunityID = uint32(intCommunityID)

	return &newPost, file, nil
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

func getLastID(r *http.Request) (uint64, error) {
	lastID := r.URL.Query().Get("id")

	if lastID == "" {
		return math.MaxInt32, nil
	}

	intLastID, err := strconv.ParseUint(lastID, 10, 32)
	if err != nil {
		return 0, err
	}

	return intLastID, nil
}

func checkAccess(userID, communityID uint32) bool {
	return true
}
