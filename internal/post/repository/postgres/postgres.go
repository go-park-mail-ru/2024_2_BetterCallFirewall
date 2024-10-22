package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

// TODO добавить сообщества
const (
	createPost      = `INSERT INTO post (author_id, content) VALUES ($1, $2) RETURNING id;`
	getPost         = `SELECT (id, author_id, content, created_at)  FROM post WHERE id = $1;`
	deletePost      = `DELETE FROM post WHERE id = $1;`
	updatePost      = `UPDATE post SET content = $1, updated_at = $2 WHERE id = $3;`
	getPostBatch    = `SELECT (id, author_id, content, created_at)  FROM post WHERE id < $1 ORDER BY created_at DESC LIMIT 10;`
	getProfilePosts = `SELECT (id, content, created_at) FROM post WHERE author_id = $1 ORDER BY created_at DESC;`
	getFriendsPost  = `SELECT (id, author_id, content, created_at) FROM post WHERE id < $1 AND author_id = ANY($2::int[]) ORDER BY created_at DESC LIMIT 10;`
	getPostAuthor   = `SELECT author_id FROM post WHERE id = $1;`
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

	if err := a.db.QueryRowContext(ctx, createPost, post.Header.AuthorID, post.PostContent.Text).Scan(&postID); err != nil {
		return 0, fmt.Errorf("postgres create post: %w", err)
	}

	return postID, nil
}

func (a *Adapter) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	var post models.Post

	if err := a.db.QueryRowContext(ctx, getPost, postID).Scan(&post.ID, &post.Header.AuthorID, &post.PostContent.Text, &post.PostContent.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrPostNotFound
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
		return myErr.ErrPostNotFound
	}

	return nil
}

func (a *Adapter) Update(ctx context.Context, post *models.Post) error {
	res, err := a.db.ExecContext(ctx, updatePost, post.PostContent.Text, post.PostContent.UpdatedAt, post.ID)

	if err != nil {
		return fmt.Errorf("postgres update post: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres update post: %w", err)
	}

	if affected == 0 {
		return myErr.ErrPostNotFound
	}

	return nil
}

func (a *Adapter) GetPosts(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	rows, err := a.db.QueryContext(ctx, getPostBatch, lastID)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get posts: %w", err)
	}

	return createPostBatchFromRows(rows)
}

func (a *Adapter) GetFriendsPosts(ctx context.Context, friendsID []uint32, lastID uint32) ([]*models.Post, error) {
	friends := convertSliceToString(friendsID)
	rows, err := a.db.QueryContext(ctx, getFriendsPost, lastID, friends)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrNoMoreContent
		}
		return nil, fmt.Errorf("postgres get friends posts: %w", err)
	}

	return createPostBatchFromRows(rows)
}

func (a *Adapter) GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	var posts []*models.Post

	rows, err := a.db.QueryContext(ctx, getProfilePosts, header.AuthorID)

	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrNoMoreContent
		}

		return nil, fmt.Errorf("postgres get author posts: %w", err)
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.PostContent.Text, &post.PostContent.CreatedAt)
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
			return 0, myErr.ErrNoMoreContent
		}
		return 0, fmt.Errorf("postgres get post author: %w", err)
	}

	return authorID, nil
}

func createPostBatchFromRows(rows *sql.Rows) ([]*models.Post, error) {
	var posts []*models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Header.AuthorID, &post.PostContent.Text, &post.PostContent.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres scan posts: %w", err)
		}
		posts = append(posts, &post)
	}

	if len(posts) < 10 {
		return posts, myErr.ErrNoMoreContent
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
