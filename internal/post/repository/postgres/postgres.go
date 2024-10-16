package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

// TODO добавить сообщества
const (
	createPost      = `INSERT INTO post (author_id, content) VALUES ($1, $2) RETURNING id;`
	getPost         = `SELECT (author_id, content, created_at)  FROM post WHERE id = $1;`
	deletePost      = `DELETE FROM post WHERE id = $1;`
	updatePost      = `UPDATE post SET content = $1, updated_at = $2 WHERE id = $3;`
	getPosts        = `SELECT (id, author_id, content, created_at)  FROM post WHERE id < $1 LIMIT 10;`
	getProfilePosts = `SELECT (id, content, created_at) FROM post WHERE author_id = $1;`
	getFriendsPost  = `SELECT (id, author_id, content, created_at) FROM post WHERE id < $1 AND author_id = ANY($2) LIMIT 10;`
	getPostAuthor   = `SELECT author_id FROM post WHERE id = $1;`
)

type Adapter struct {
	db *pgxpool.Conn
}

func NewAdapter(db *pgxpool.Conn) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) Create(ctx context.Context, post *models.Post) (uint32, error) {
	var postID uint32

	if err := a.db.QueryRow(ctx, createPost, post.Header.AuthorID, post.PostContent.Text).Scan(&postID); err != nil {
		return 0, fmt.Errorf("postgres create post: %w", err)
	}

	return postID, nil
}

func (a *Adapter) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	var post models.Post

	if err := a.db.QueryRow(ctx, getPost, postID).Scan(&post.Header.AuthorID, &post.PostContent.Text, &post.PostContent.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myErr.ErrPostNotFound
		}

		return nil, fmt.Errorf("postgres get post: %w", err)
	}
	post.ID = postID

	return &post, nil
}

func (a *Adapter) Delete(ctx context.Context, postID uint32) error {
	_, err := a.db.Exec(ctx, deletePost, postID)
	if errors.Is(err, pgx.ErrNoRows) {
		return myErr.ErrPostNotFound
	}

	if err != nil {
		return fmt.Errorf("postgres delete post: %w", err)
	}

	return nil
}

func (a *Adapter) Update(ctx context.Context, post *models.Post) error {
	_, err := a.db.Exec(ctx, updatePost, post.PostContent.Text, post.PostContent.UpdatedAt, post.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return myErr.ErrPostNotFound
	}

	if err != nil {
		return fmt.Errorf("postgres update post: %w", err)
	}

	return nil
}

func (a *Adapter) GetPosts(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	row, err := a.db.Query(ctx, getPosts, lastID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, myErr.ErrNoMoreContent
	}

	if err != nil {
		return nil, fmt.Errorf("postgres get posts: %w", err)
	}
	defer row.Close()

	return createPostBatchFromRows(row)
}

func (a *Adapter) GetFriendsPosts(ctx context.Context, friendsID []uint32, lastID uint32) ([]*models.Post, error) {
	rows, err := a.db.Query(ctx, getFriendsPost, lastID, friendsID)
	if err != nil {
		return nil, fmt.Errorf("postgres get friends posts: %w", err)
	}
	defer rows.Close()

	return createPostBatchFromRows(rows)
}

func (a *Adapter) GetAuthorsPosts(ctx context.Context, authorID uint32) ([]*models.Post, error) {
	var (
		post  models.Post
		posts []*models.Post
	)

	rows, err := a.db.Query(ctx, getProfilePosts, authorID)
	if err != nil {
		return nil, fmt.Errorf("postgres get author posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(post.ID, post.PostContent.Text, post.PostContent.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("postgres get author posts: %w", err)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (a *Adapter) GetPostAuthor(ctx context.Context, postID uint32) (uint32, error) {
	var authorID uint32

	if err := a.db.QueryRow(ctx, getPostAuthor, postID).Scan(&authorID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, myErr.ErrPostNotFound
		}

		return 0, fmt.Errorf("postgres get post author: %w", err)
	}

	return authorID, nil
}

func createPostBatchFromRows(rows pgx.Rows) ([]*models.Post, error) {
	var (
		post  models.Post
		posts []*models.Post
	)

	for rows.Next() {
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
