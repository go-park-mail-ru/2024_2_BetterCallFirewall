package postgres

import (
	"database/sql"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/post/entities"
)

// TODO добавить сообщества
const (
	createPostTable = `CREATE TABLE IF NOT EXISTS post (id SERIAL PRIMARY KEY, author_id INTEGER REFERENCES profile(id) ON DELETE CASCADE, content_id INTEGER REFERENCES content(id) ON DELETE CASCADE);`
	createPost      = `INSERT INTO post (author_id, content_id) VALUES ($1, $2);`
	getPost         = `SELECT (author_id, content_id, text, created_at, image_path)  FROM post AS p INNER JOIN content AS c ON c.id = p.content_id INNER JOIN content_image AS ci ON ci.content_id = p.content_id WHERE id = $1;`
	deletePost      = `DELETE FROM post WHERE id = $1;`
)

type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) CreateNewTable() error {
	_, err := a.db.Exec(createPostTable)
	if err != nil {
		return fmt.Errorf("create post table: %w", err)
	}

	return nil
}

func (a *Adapter) Create(post *entities.PostDB) (uint32, error) {
	res, err := a.db.Exec(createPost, post.AuthorID, post.ContentID)
	if err != nil {
		return 0, fmt.Errorf("postgres create post: %w", err)
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("postgres create post: %w", err)
	}

	return uint32(postID), nil
}

func (a *Adapter) Get(postID uint32) (*models.Post, error) {
	var (
		post      models.Post
		contentID uint32
	)

	row := a.db.QueryRow(getPost, postID)

	err := row.Scan(&post.AuthorID, &contentID, &post.PostContent.Text, &post.PostContent.CreatedAt, &post.PostContent.File)
	if err != nil {
		return nil, fmt.Errorf("postgres get post: %w", err)
	}

	return &post, nil
}

func (a *Adapter) Delete(postID uint32) error {
	_, err := a.db.Exec(deletePost, postID)
	if err != nil {
		return fmt.Errorf("postgres delete post: %w", err)
	}

	return nil
}
