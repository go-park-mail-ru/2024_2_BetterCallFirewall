package postgres

import (
	"database/sql"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

const (
	createTable           = `CREATE TABLE IF NOT EXISTS content (id SERIAL PRIMARY KEY, text TEXT, created_at DATE NOT NULL, updated_at DATE NOT NULL);`
	createTableManyToMany = `CREATE TABLE IF NOT EXISTS content_image (content_id INTEGER REFERENCES content(id) ON DELETE CASCADE, image_path TEXT NOT NULL);`
	createContent         = `INSERT INTO content (text, created_at, update_at) VALUES ($1, $2, $3);`
	createContentImage    = `INSERT INTO content_image (content_id, image_path) VALUES ($1, $2);`
	updateContent         = `UPDATE content SET text = $1, updated_at = $2 WHERE id = $3;`
	updateImage           = `UPDATE content_image SET image_path = $1 WHERE id = $2;`
)

type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) Create(content *models.Content) (uint32, error) {
	res, err := a.db.Exec(createContent, content.Text, content.CreatedAt, content.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("psql creating content: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("psql get last inserted id: %w", err)
	}

	res, err = a.db.Exec(createContentImage, id, content.File.Path)

	return uint32(id), nil
}

func (a *Adapter) Update(content *models.Content) error {
	_, err := a.db.Exec(updateContent, content.Text, content.UpdatedAt, content.ID)
	if err != nil {
		return fmt.Errorf("psql updating content: %w", err)
	}

	_, err = a.db.Exec(updateImage, content.File.Path, content.ID)
	if err != nil {
		return fmt.Errorf("psql updating image: %w", err)
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
