package postgres

import (
	"database/sql"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

// TODO добавить сообщества
const (
	createPostTable = `CREATE TABLE IF NOT EXISTS post (id SERIAL PRIMARY KEY, author_id INTEGER REFERENCES profile(id) ON DELETE CASCADE, content_id INTEGER REFERENCES content(id) ON DELETE CASCADE);`
	createPost      = `INSERT INTO post (author_id, content_id) VALUES ($1, $2);`
	getPost         = `SELECT (author_id, content_id) FROM post WHERE id = $1;`
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
		return err
	}
	return nil
}

func (a *Adapter) Create(post *models.PostDB) (uint32, error) {

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

func (a *Adapter) Get(postID uint32) (*models.PostDB, error) {
	var post models.PostDB
	row := a.db.QueryRow(getPost, postID)

	err := row.Scan(&post.AuthorID, &post.ContentID)
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
