package postgres

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/post/entities"
)

// TODO добавить сообщества
const (
	createPostTable = `CREATE TABLE IF NOT EXISTS post (id INT PRIMARY KEY, author_id INTEGER REFERENCES profile(id) ON DELETE CASCADE, content TEXT, created_at DATE, updated_at DATE);`
	createPost      = `INSERT INTO post (id, author_id, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);`
	getPost         = `SELECT (author_id, content, created_at)  FROM post WHERE id = $1;`
	deletePost      = `DELETE FROM post WHERE id = $1;`
	checkCreater    = `SELECT author_id FROM post WHERE id = $1;`
	getPosts        = `SELECT (id, author_id, content, created_at)  FROM post WHERE id < $1 LIMIT 10;`
	getProfilePosts = `SELECT (id, content, created_at) FROM post WHERE author_id = $1;`
)

type Adapter struct {
	db      *sql.DB
	counter uint32
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db:      db,
		counter: 1,
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
	_, err := a.db.Exec(createPost, a.counter, post.AuthorID, post.Content, post.Created, post.Updated)
	if err != nil {
		return 0, fmt.Errorf("postgres create post: %w", err)
	}
	a.counter++

	return a.counter - 1, nil
}

func (a *Adapter) Get(postID uint32) (*models.Post, error) {
	var post models.Post

	row := a.db.QueryRow(getPost, postID)

	err := row.Scan(&post.AuthorID, &post.PostContent.Text, &post.PostContent.CreatedAt)
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

func (a *Adapter) Update(post *entities.PostDB) error {
	return nil
}

func (a *Adapter) CheckAccess(profileID uint32, postID uint32) (bool, error) {
	row, err := a.db.Query(checkCreater, postID)
	if err != nil {
		return false, fmt.Errorf("postgres check access: %w", err)
	}
	defer row.Close()

	var createrID uint32
	err = row.Scan(&createrID)
	if err != nil {
		return false, fmt.Errorf("postgres check access: %w", err)
	}

	if createrID != profileID {
		return false, nil
	}

	return true, nil
}

func (a *Adapter) GetPosts(lastID uint32, newRequest bool) ([]*models.Post, error) {
	if newRequest {
		lastID = a.counter
	}

	row, err := a.db.Query(getPosts, lastID)
	if err != nil {
		return nil, fmt.Errorf("postgres get posts: %w", err)
	}
	defer row.Close()

	return createPostBatchFromRows(row)
}

func (a *Adapter) GetFriendsPosts(friendsID []uint32, lastID uint32, newRequest bool) ([]*models.Post, error) {
	if newRequest {
		lastID = a.counter
	}

	var paramrefs string

	for i := range friendsID {
		paramrefs += `$` + strconv.Itoa(i+2) + `,`
	}
	paramrefs = paramrefs[:len(paramrefs)-1]

	query := `SELECT (id, author_id, content, created_at) FROM post
	WHERE id < $1 AND author_id IN (` + paramrefs + `)LIMIT 10;`

	rows, err := a.db.Query(query, lastID, friendsID)
	if err != nil {
		return nil, fmt.Errorf("postgres get friends posts: %w", err)
	}
	defer rows.Close()

	return createPostBatchFromRows(rows)
}

func (a *Adapter) GetAuthorsPosts(authorID uint32) ([]*models.Post, error) {
	var (
		post  models.Post
		posts []*models.Post
	)

	rows, err := a.db.Query(getProfilePosts, authorID)
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

func createPostBatchFromRows(rows *sql.Rows) ([]*models.Post, error) {
	var (
		post  models.Post
		posts []*models.Post
	)

	for rows.Next() {
		err := rows.Scan(&post.ID, &post.AuthorID, &post.PostContent.Text, &post.PostContent.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("postgres scan posts: %w", err)
		}
		posts = append(posts, &post)
	}

	if len(posts) < 10 {
		return posts, myErr.ErrNoMoreContent
	}

	return posts, nil
}
