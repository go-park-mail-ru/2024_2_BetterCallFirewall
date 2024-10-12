package postgres

import (
	"database/sql"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

const (
	createTable           = `CREATE TABLE IF NOT EXISTS content (id INT PRIMARY KEY, text TEXT, created_at DATE NOT NULL, updated_at DATE NOT NULL);`
	createTableManyToMany = `CREATE TABLE IF NOT EXISTS content_image (content_id INTEGER REFERENCES content(id) ON DELETE CASCADE, image_path TEXT NOT NULL);`
	createContent         = `INSERT INTO content (id, text, created_at, update_at) VALUES ($1, $2, $3, $4);`
	createContentImage    = `INSERT INTO content_image (content_id, image_path) VALUES ($1, $2);`
	updateContent         = `UPDATE content SET text = $1, updated_at = $2 WHERE id = $3;`
	updateImage           = `UPDATE content_image SET image_path = $1 WHERE id = $2;`
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

func (a *Adapter) Create(content *models.Content) (uint32, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("start transaction: %w", err)
	}

	_, err = tx.Exec(createContent, a.counter, content.Text, content.CreatedAt, content.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("psql creating content: %w", err)
	}

	_, err = tx.Exec(createContentImage, a.counter, content.File.Path)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("psql creating content_image: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}
	a.counter++

	return a.counter - 1, nil
}

func (a *Adapter) Update(content *models.Content) error {
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	_, err = tx.Exec(updateContent, content.Text, content.UpdatedAt, content.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("psql updating content: %w", err)
	}

	_, err = tx.Exec(updateImage, content.File.Path, content.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("psql updating image: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (a *Adapter) NewTable() error {
	_, err := a.db.Exec(createTable)
	if err != nil {
		return fmt.Errorf("create table content: %w", err)
	}

	_, err = a.db.Exec(createTableManyToMany)
	if err != nil {
		return fmt.Errorf("create table content_image: %w", err)
	}

	return nil
}
