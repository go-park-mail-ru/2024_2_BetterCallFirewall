package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

const ()

type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	return &Adapter{
		db: db,
	}
}

func (a *Adapter) Create(user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (email) DO NOTHING;
`
	res, err := a.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("postgres create user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres create user rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("postgres create user: %w", myErr.ErrUserAlreadyExists)
	}

	return nil
}

func (a *Adapter) GetByEmail(email string) (*models.User, error) {
	query := `SELECT first_name, last_name, email, password FROM users WHERE email = $1`
	row := a.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres get user: %w", myErr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("postgres get user: %w", err)
	}

	return &user, nil
}

func (a *Adapter) CreateNewUserTable() error {
	newTableString := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE ,
		password VARCHAR(100) NOT NULL
	);`

	_, err := a.db.Exec(newTableString)
	if err != nil {
		return fmt.Errorf("postgres create user table: %w", err)
	}

	return nil
}

func (a *Adapter) CreateNewSessionTable() error {

}

func (a *Adapter) CreateSession(sess *models.Session) error {
	//TODO
	return nil
}

func (a *Adapter) FindSession(sessID string) (*models.Session, error) {
	//TODO
	return nil, nil
}

func (a *Adapter) DeleteSession(sessID string) error {
	return nil
}

func StartPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	return db, nil
}
