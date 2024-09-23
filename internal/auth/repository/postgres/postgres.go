package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

const (
	CreateNewSessionTable = `CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    sess_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    
)`
	CreateSession = `INSERT INTO sessions (sess_id, user_id) VALUES ($1, $2)`
	FindSession   = `SELECT sess_id, user_id FROM sessions WHERE sess_id = $1`
	DeleteSession = `DELETE FROM sessions WHERE sess_id = $1`
)

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
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE ,
		password TEXT NOT NULL
	);`

	_, err := a.db.Exec(newTableString)
	if err != nil {
		return fmt.Errorf("postgres create user table: %w", err)
	}

	return nil
}

func (a *Adapter) CreateNewSessionTable() error {
	_, err := a.db.Exec(CreateNewSessionTable)
	if err != nil {
		return fmt.Errorf("postgres create session table: %w", err)
	}
	return nil
}

func (a *Adapter) CreateSession(sess *models.Session) error {
	res, err := a.db.Exec(CreateSession, sess.ID, sess.UserID)
	if err != nil {
		return fmt.Errorf("postgres create session table: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres create session rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("postgres create session: %w", myErr.ErrSessionNotFound)
	}

	return nil
}

func (a *Adapter) FindSession(sessID string) (*models.Session, error) {
	res := a.db.QueryRow(FindSession, sessID)
	var sess models.Session
	err := res.Scan(&sess.ID, &sess.UserID)
	if err != nil {
		return nil, fmt.Errorf("postgres find session table: %w", err)
	}

	return &sess, nil
}

func (a *Adapter) DeleteSession(sessID string) error {
	_, err := a.db.Exec(DeleteSession, sessID)
	if err != nil {
		return fmt.Errorf("postgres delete session table: %w", err)
	}
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
