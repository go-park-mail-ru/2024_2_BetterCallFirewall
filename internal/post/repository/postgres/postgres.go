package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const (
	createPost      = `INSERT INTO post (author_id, content, file_path) VALUES ($1, $2, $3) RETURNING id;`
	getPost         = `SELECT id, author_id, content, file_path, created_at  FROM post WHERE id = $1;`
	deletePost      = `DELETE FROM post WHERE id = $1;`
	updatePost      = `UPDATE post SET content = $1, updated_at = $2, file_path = $3 WHERE id = $4;`
	getPostBatch    = `SELECT id, CASE WHEN author_id IS NULL THEN 0 ELSE author_id END, CASE WHEN community_id IS NULL THEN 0 ELSE community_id END, content, file_path, created_at  FROM post WHERE id < $1 ORDER BY created_at DESC LIMIT 10;`
	getProfilePosts = `SELECT id, content, file_path, created_at FROM post WHERE author_id = $1 ORDER BY created_at DESC;`
	getFriendsPost  = `SELECT id, author_id, content, file_path, created_at FROM post WHERE id < $1 AND author_id = ANY($2::int[]) ORDER BY created_at DESC LIMIT 10;`
	getPostAuthor   = `SELECT author_id FROM post WHERE id = $1;`

	createCommunityPost = `INSERT INTO post (community_id, content, file_path) VALUES ($1, $2, $3) RETURNING id;`
	getCommunityPosts   = `SELECT id, community_id, content, file_path, created_at FROM post WHERE community_id = $1 AND id < $2 ORDER BY id DESC LIMIT 10;`

	AddLikeToPost      = `INSERT INTO reaction (post_id, user_id) VALUES ($1, $2);`
	DeleteLikeFromPost = `DELETE FROM reaction WHERE post_id = $1 AND user_id = $2;`
	GetLikesOnPost     = `SELECT COUNT(*) FROM reaction WHERE post_id = $1;`
	CheckLike          = `SELECT COUNT(*) FROM reaction WHERE post_id = $1 AND user_id=$2;`

	createComment      = `INSERT INTO comment (user_id, post_id, content, file_path) VALUES ($1, $2, $3, $4) RETURNING id;`
	updateComment      = `UPDATE comment SET content = $1, file_path = $2, updated_at = NOW() WHERE id = $3;`
	deleteComment      = `DELETE FROM comment WHERE id = $1;`
	getCommentsBatch   = `SELECT id, user_id, content, file_path, created_at FROM comment WHERE post_id = $1 and id < $2 ORDER BY created_at DESC LIMIT 10;`
	getCommentBatchAsc = `SELECT id, user_id, content, file_path, created_at FROM comment WHERE post_id = $1 and id > $2 order by created_at LIMIT 10;`
	getCommentAuthor   = `SELECT user_id FROM comment WHERE id = $1`
	getCommentCount    = `SELECT COUNT(*) FROM comment WHERE post_id=$1`
)

type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) Create(ctx context.Context, post *models.Post) (uint32, error) {
	var postID uint32

	if err := a.db.QueryRowContext(
		ctx, createPost, post.Header.AuthorID, post.PostContent.Text, post.PostContent.File,
	).Scan(&postID); err != nil {
		return 0, fmt.Errorf("postgres create post: %w", err)
	}

	return postID, nil
}

func (a *Adapter) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	var post models.Post

	if err := a.db.QueryRowContext(ctx, getPost, postID).
		Scan(
			&post.ID, &post.Header.AuthorID, &post.PostContent.Text, &post.PostContent.File,
			&post.PostContent.CreatedAt,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrPostNotFound
		}

		return nil, fmt.Errorf("postgres get post: %w", err)
	}

	return &post, nil
}

func (a *Adapter) Delete(ctx context.Context, postID uint32) error {
	res, err := a.db.ExecContext(ctx, deletePost, postID)

	if err != nil {
		return fmt.Errorf("postgres delete post: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres delete post: %w", err)
	}

	if affected == 0 {
		return my_err.ErrPostNotFound
	}

	return nil
}

func (a *Adapter) Update(ctx context.Context, post *models.Post) error {
	res, err := a.db.ExecContext(
		ctx, updatePost, post.PostContent.Text, post.PostContent.UpdatedAt, post.PostContent.File, post.ID,
	)

	if err != nil {
		return fmt.Errorf("postgres update post: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres update post: %w", err)
	}

	if affected == 0 {
		return my_err.ErrPostNotFound
	}

	return nil
}

func (a *Adapter) GetPosts(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	rows, err := a.db.QueryContext(ctx, getPostBatch, lastID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID, &post.Header.AuthorID, &post.Header.CommunityID,
			&post.PostContent.Text, &post.PostContent.File, &post.PostContent.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres scan posts: %w", err)
		}
		posts = append(posts, &post)
	}

	if len(posts) == 0 {
		return posts, my_err.ErrNoMoreContent
	}

	return posts, nil
}

func (a *Adapter) GetFriendsPosts(ctx context.Context, friendsID []uint32, lastID uint32) ([]*models.Post, error) {
	friends := convertSliceToString(friendsID)
	rows, err := a.db.QueryContext(ctx, getFriendsPost, lastID, friends)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get friends posts: %w", err)
	}

	return createPostBatchFromRows(rows)
}

func (a *Adapter) GetAuthorPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	var posts []*models.Post

	rows, err := a.db.QueryContext(ctx, getProfilePosts, header.AuthorID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}

		return nil, fmt.Errorf("postgres get author posts: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.PostContent.Text, &post.PostContent.File, &post.PostContent.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("postgres get author posts: %w", err)
		}
		post.Header = *header
		posts = append(posts, &post)
	}

	return posts, nil
}

func (a *Adapter) GetPostAuthor(ctx context.Context, postID uint32) (uint32, error) {
	var authorID uint32

	if err := a.db.QueryRowContext(ctx, getPostAuthor, postID).Scan(&authorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, my_err.ErrPostNotFound
		}
		return 0, fmt.Errorf("postgres get post author: %w", err)
	}

	return authorID, nil
}

func createPostBatchFromRows(rows *sql.Rows) ([]*models.Post, error) {
	var posts []*models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(
			&post.ID, &post.Header.AuthorID, &post.PostContent.Text, &post.PostContent.File,
			&post.PostContent.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres scan posts: %w", err)
		}
		posts = append(posts, &post)
	}

	if len(posts) == 0 {
		return posts, my_err.ErrNoMoreContent
	}

	return posts, nil
}

func convertSliceToString(sl []uint32) string {
	var sb strings.Builder
	sb.Grow(len(sl) * 3)

	sb.WriteString("{")
	for _, v := range sl {
		sb.WriteString(fmt.Sprintf("%d, ", v))
	}
	res := sb.String()
	res = strings.TrimSuffix(res, ", ")
	res += "}"

	return res
}

func (a *Adapter) CreateCommunityPost(ctx context.Context, post *models.Post, communityID uint32) (uint32, error) {
	var ID uint32
	if err := a.db.QueryRowContext(
		ctx, createCommunityPost, communityID, post.PostContent.Text, post.PostContent.File,
	).Scan(&ID); err != nil {
		return 0, fmt.Errorf("postgres create community post db: %w", err)
	}

	return ID, nil

}

func (a *Adapter) GetCommunityPosts(ctx context.Context, communityID, id uint32) ([]*models.Post, error) {
	var posts []*models.Post
	rows, err := a.db.QueryContext(ctx, getCommunityPosts, communityID, id)
	if err != nil {
		return nil, fmt.Errorf("postgres get community posts: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		err = rows.Scan(
			&post.ID, &post.Header.CommunityID, &post.PostContent.Text, &post.PostContent.File,
			&post.PostContent.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres get community posts: %w", err)
		}
		posts = append(posts, post)
	}
	if len(posts) == 0 {
		return posts, my_err.ErrNoMoreContent
	}

	return posts, nil
}

func (a *Adapter) SetLikeToPost(ctx context.Context, postID uint32, userID uint32) error {
	res, err := a.db.ExecContext(ctx, AddLikeToPost, postID, userID)
	if num, err := res.RowsAffected(); err == nil && num == 0 {
		return my_err.ErrLikeAlreadyExists
	}
	if err != nil {
		return err
	}
	return nil
}

func (a *Adapter) DeleteLikeFromPost(ctx context.Context, postID uint32, userID uint32) error {
	_, err := a.db.ExecContext(ctx, DeleteLikeFromPost, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (a *Adapter) GetLikesOnPost(ctx context.Context, postID uint32) (uint32, error) {
	var likes uint32
	err := a.db.QueryRowContext(ctx, GetLikesOnPost, postID).Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, my_err.ErrWrongPost
		}
		return 0, err
	}
	return likes, nil
}

func (a *Adapter) CheckLikes(ctx context.Context, postID, userID uint32) (bool, error) {
	var likes uint32
	err := a.db.QueryRowContext(ctx, CheckLike, postID, userID).Scan(&likes)
	if err != nil {
		return false, fmt.Errorf("postgres check likes on post %d: %w", postID, err)
	}

	if likes == 0 {
		return false, nil
	}

	return true, nil
}

func (a *Adapter) CreateComment(ctx context.Context, comment *models.Content, userID, postID uint32) (uint32, error) {
	var id uint32

	if err := a.db.QueryRowContext(
		ctx, createComment, userID, postID, comment.Text, comment.File,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("postgres create comment: %w", err)
	}

	return id, nil
}

func (a *Adapter) DeleteComment(ctx context.Context, commentID uint32) error {
	res, err := a.db.ExecContext(ctx, deleteComment, commentID)

	if err != nil {
		return fmt.Errorf("postgres delete comment: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres delete comment: %w", err)
	}

	if affected == 0 {
		return my_err.ErrWrongComment
	}

	return nil
}

func (a *Adapter) UpdateComment(ctx context.Context, comment *models.Content, commentID uint32) error {
	res, err := a.db.ExecContext(
		ctx, updateComment, comment.Text, comment.File, commentID,
	)

	if err != nil {
		return fmt.Errorf("postgres update comment: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres update comment: %w", err)
	}

	if affected == 0 {
		return my_err.ErrWrongComment
	}

	return nil
}

func (a *Adapter) GetComments(ctx context.Context, postID, lastID uint32, newest bool) ([]*models.Comment, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if newest {
		rows, err = a.db.QueryContext(ctx, getCommentsBatch, postID, lastID)
	} else {
		rows, err = a.db.QueryContext(ctx, getCommentBatchAsc, postID, lastID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get posts: %w", err)
	}

	var comments []*models.Comment
	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(
			&comment.ID, &comment.Header.AuthorID, &comment.Content.Text, &comment.Content.File,
			&comment.Content.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres get comments: %w", err)
		}

		comments = append(comments, &comment)
	}

	if len(comments) == 0 {
		return comments, my_err.ErrNoMoreContent
	}

	return comments, nil
}

func (a *Adapter) GetCommentAuthor(ctx context.Context, commentID uint32) (uint32, error) {
	var authorID uint32

	if err := a.db.QueryRowContext(ctx, getCommentAuthor, commentID).Scan(&authorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, my_err.ErrWrongComment
		}
		return 0, fmt.Errorf("postgres get comment author: %w", err)
	}

	return authorID, nil
}

func (a *Adapter) GetCommentCount(ctx context.Context, postID uint32) (uint32, error) {
	var count uint32

	if err := a.db.QueryRowContext(ctx, getCommentCount, postID).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, my_err.ErrWrongPost
		}
		return 0, fmt.Errorf("postgres get comment count: %w", err)
	}

	return count, nil
}
