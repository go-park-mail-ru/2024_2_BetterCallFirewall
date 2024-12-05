package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const (
	postIDkey    = "id"
	commentIDKey = "comment_id"
	filePrefix   = "/image/"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type PostService interface {
	Create(ctx context.Context, post *models.Post) (uint32, error)
	Get(ctx context.Context, postID, userID uint32) (*models.Post, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, postID uint32) error
	GetBatch(ctx context.Context, lastID, userID uint32) ([]*models.Post, error)
	GetBatchFromFriend(ctx context.Context, userID uint32, lastID uint32) ([]*models.Post, error)
	GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error)

	GetCommunityPost(ctx context.Context, communityID, userID, lastID uint32) ([]*models.Post, error)
	CreateCommunityPost(ctx context.Context, post *models.Post) (uint32, error)
	CheckAccessToCommunity(ctx context.Context, userID uint32, communityID uint32) bool

	SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error
	DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error
	CheckLikes(ctx context.Context, postID, userID uint32) (bool, error)
}

type CommentService interface {
	Comment(ctx context.Context, userID, postID uint32, comment *models.Content) (*models.Comment, error)
	DeleteComment(ctx context.Context, commentID, userID uint32) error
	EditComment(ctx context.Context, commentID, userID uint32, comment *models.Content) error
	GetComments(ctx context.Context, postID, lastID uint32) ([]*models.Comment, error)
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestId string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestId string)

	ErrorInternal(w http.ResponseWriter, err error, requestId string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestId string)
	LogError(err error, requestId string)
}

type PostController struct {
	postService    PostService
	commentService CommentService
	responder      Responder
}

func NewPostController(service PostService, commentService CommentService, responder Responder) *PostController {
	return &PostController{
		postService:    service,
		commentService: commentService,
		responder:      responder,
	}
}

func (pc *PostController) Create(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		comunity  = r.URL.Query().Get("community")
		id        uint32
		err       error
	)

	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	newPost, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !validateContent(newPost.PostContent) {
		pc.responder.ErrorBadRequest(w, my_err.ErrTextTooLong, reqID)
		return
	}
	if comunity != "" {
		comID, err := strconv.ParseUint(comunity, 10, 32)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if !pc.checkAccessToCommunity(r, uint32(comID)) {
			pc.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
			return
		}

		newPost.Header.CommunityID = uint32(comID)
		id, err = pc.postService.CreateCommunityPost(r.Context(), newPost)
		if err != nil {
			pc.responder.ErrorInternal(w, err, reqID)
			return
		}
	} else {
		id, err = pc.postService.Create(r.Context(), newPost)
		if err != nil {
			pc.responder.ErrorInternal(w, fmt.Errorf("create controller: %w", err), reqID)
			return
		}
	}

	newPost.ID = id

	pc.responder.OutputJSON(w, newPost, reqID)
}

func (pc *PostController) GetOne(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	postID, err := getIDFromURL(r, postIDkey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	post, err := pc.postService.Get(r.Context(), postID, sess.UserID)
	if err != nil {
		if errors.Is(err, my_err.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if !errors.Is(err, my_err.ErrAnotherService) {
			pc.responder.ErrorInternal(w, err, reqID)
			return
		}
	}

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) Update(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = getIDFromURL(r, postIDkey)
		community = r.URL.Query().Get("community")
	)

	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	if community != "" {
		comID, err := strconv.ParseUint(community, 10, 32)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if !pc.checkAccessToCommunity(r, uint32(comID)) {
			pc.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
			return
		}
	} else {
		if !pc.checkAccess(r, id) {
			pc.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
			return
		}
	}

	post, err := pc.getPostFromBody(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !validateContent(post.PostContent) {
		pc.responder.ErrorBadRequest(w, my_err.ErrTextTooLong, reqID)
		return
	}
	post.ID = id

	if err := pc.postService.Update(r.Context(), post); err != nil {
		if errors.Is(err, my_err.ErrPostNotFound) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, post, reqID)
}

func (pc *PostController) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok   = r.Context().Value("requestID").(string)
		postID, err = getIDFromURL(r, postIDkey)
		community   = r.URL.Query().Get("community")
	)

	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if community != "" {
		comID, err := strconv.ParseUint(community, 10, 32)
		if err != nil {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if !pc.checkAccessToCommunity(r, uint32(comID)) {
			pc.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
			return
		}
	} else {
		if !pc.checkAccess(r, postID) {
			pc.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
			return
		}
	}

	if err := pc.postService.Delete(r.Context(), postID); err != nil {
		if errors.Is(err, my_err.ErrPostNotFound) {
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
		reqID, ok   = r.Context().Value("requestID").(string)
		section     = r.URL.Query().Get("section")
		communityID = r.URL.Query().Get("community")
		posts       []*models.Post
		intLastID   uint64
		err         error
		id          uint64
	)

	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	intLastID, err = getLastID(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	switch section {
	case "friend":
		{
			posts, err = pc.postService.GetBatchFromFriend(r.Context(), sess.UserID, uint32(intLastID))
		}
	case "":
		{
			if communityID != "" {
				id, err = strconv.ParseUint(communityID, 10, 32)
				if err != nil {
					pc.responder.ErrorBadRequest(w, err, reqID)
					return
				}
				posts, err = pc.postService.GetCommunityPost(r.Context(), uint32(id), sess.UserID, uint32(intLastID))
			} else {
				posts, err = pc.postService.GetBatch(r.Context(), uint32(intLastID), sess.UserID)
			}
		}
	default:
		pc.responder.ErrorBadRequest(w, my_err.ErrInvalidQuery, reqID)
		return
	}

	if err != nil {
		if errors.Is(err, my_err.ErrNoMoreContent) {
			pc.responder.OutputNoMoreContentJSON(w, reqID)
			return
		}
		if !errors.Is(err, my_err.ErrAnotherService) {
			pc.responder.ErrorInternal(w, err, reqID)
			return
		}
	}

	pc.responder.OutputJSON(w, posts, reqID)
}

func (pc *PostController) getPostFromBody(r *http.Request) (*models.Post, error) {
	var newPost models.Post

	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		return nil, err
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return nil, err
	}
	newPost.Header.AuthorID = sess.UserID

	return &newPost, nil
}

func getIDFromURL(r *http.Request, key string) (uint32, error) {
	vars := mux.Vars(r)

	id := vars[key]
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

func (pc *PostController) checkAccess(r *http.Request, postID uint32) bool {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return false
	}

	userID := sess.UserID
	authorID, err := pc.postService.GetPostAuthorID(r.Context(), postID)
	if err != nil {
		return false
	}

	if userID != authorID {
		return false
	}

	return true
}

func (pc *PostController) checkAccessToCommunity(r *http.Request, communityID uint32) bool {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return false
	}
	userID := sess.UserID

	return pc.postService.CheckAccessToCommunity(r.Context(), userID, communityID)
}

func (pc *PostController) SetLikeOnPost(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	postID, err := getIDFromURL(r, postIDkey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, errSession := models.SessionFromContext(r.Context())
	if errSession != nil {
		pc.responder.ErrorBadRequest(w, errSession, reqID)
		return
	}

	set, err := pc.postService.CheckLikes(r.Context(), postID, sess.UserID)
	if err != nil {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if set {
		pc.responder.ErrorBadRequest(w, my_err.ErrInvalidQuery, reqID)
		return
	}

	err = pc.postService.SetLikeToPost(r.Context(), postID, sess.UserID)
	if err != nil {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, "like is set on post", reqID)
}

func (pc *PostController) DeleteLikeFromPost(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	postID, err := getIDFromURL(r, postIDkey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	sess, errSession := models.SessionFromContext(r.Context())
	if errSession != nil {
		pc.responder.ErrorBadRequest(w, errSession, reqID)
		return
	}

	set, err := pc.postService.CheckLikes(r.Context(), postID, sess.UserID)
	if err != nil {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	if !set {
		pc.responder.ErrorBadRequest(w, my_err.ErrInvalidQuery, reqID)
		return
	}

	err = pc.postService.DeleteLikeFromPost(r.Context(), postID, sess.UserID)
	if err != nil {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, "like is unset from post", reqID)
}

func (pc *PostController) Comment(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	postID, err := getIDFromURL(r, postIDkey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, errSession := models.SessionFromContext(r.Context())
	if errSession != nil {
		pc.responder.ErrorBadRequest(w, errSession, reqID)
		return
	}

	var content models.Content
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !validateContent(content) {
		pc.responder.ErrorBadRequest(w, my_err.ErrTextTooLong, reqID)
		return
	}

	newComment, err := pc.commentService.Comment(r.Context(), sess.UserID, postID, &content)
	if err != nil {
		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, newComment, reqID)
}

func (pc *PostController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	commentID, err := getIDFromURL(r, commentIDKey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, errSession := models.SessionFromContext(r.Context())
	if errSession != nil {
		pc.responder.ErrorBadRequest(w, errSession, reqID)
		return
	}

	err = pc.commentService.DeleteComment(r.Context(), commentID, sess.UserID)
	if err != nil {
		if errors.Is(err, my_err.ErrAccessDenied) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}

		if errors.Is(err, my_err.ErrWrongComment) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}

		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, "comment is deleted", reqID)
}

func (pc *PostController) EditComment(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	commentID, err := getIDFromURL(r, commentIDKey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, errSession := models.SessionFromContext(r.Context())
	if errSession != nil {
		pc.responder.ErrorBadRequest(w, errSession, reqID)
		return
	}

	var content models.Content
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !validateContent(content) {
		pc.responder.ErrorBadRequest(w, my_err.ErrTextTooLong, reqID)
		return
	}

	if err := pc.commentService.EditComment(r.Context(), commentID, sess.UserID, &content); err != nil {
		if errors.Is(err, my_err.ErrAccessDenied) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}

		if errors.Is(err, my_err.ErrWrongComment) {
			pc.responder.ErrorBadRequest(w, err, reqID)
			return
		}

		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, "comment is updated", reqID)
}

func (pc *PostController) GetComments(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		pc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	postID, err := getIDFromURL(r, postIDkey)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	lastId, err := getLastID(r)
	if err != nil {
		pc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	comments, err := pc.commentService.GetComments(r.Context(), postID, uint32(lastId))
	if err != nil {
		if errors.Is(err, my_err.ErrNoMoreContent) {
			pc.responder.OutputNoMoreContentJSON(w, reqID)
			return
		}

		pc.responder.ErrorInternal(w, err, reqID)
		return
	}

	pc.responder.OutputJSON(w, comments, reqID)
}

func validateContent(content models.Content) bool {
	return validateFile(content.File) && len(content.Text) < 500
}

func validateFile(filepath models.Picture) bool {
	return len(filepath) < 100 && (len(filepath) == 0 || strings.HasPrefix(string(filepath), filePrefix))
}
