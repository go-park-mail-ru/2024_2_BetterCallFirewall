package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

const (
	CreateUserTable       = `CREATE TABLE IF NOT EXISTS person (id SERIAL PRIMARY KEY, first_name TEXT NOT NULL CONSTRAINT first_name_length CHECK (CHAR_LENGTH(first_name) <= 30), last_name TEXT NOT NULL CONSTRAINT last_name_length CHECK (CHAR_LENGTH(last_name) <= 30), email TEXT NOT NULL UNIQUE NOT NULL CONSTRAINT email_length CHECK (CHAR_LENGTH(email) <= 50), password TEXT NOT NULL CONSTRAINT password_length CHECK (CHAR_LENGTH(password) <= 50));`
	CreateUser            = `INSERT INTO person (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING;`
	GetUserByEmail        = `SELECT id, first_name, last_name, email, password FROM person WHERE email = $1;`
	CreateNewSessionTable = `CREATE TABLE IF NOT EXISTS session (id SERIAL PRIMARY KEY, sess_id TEXT NOT NULL, user_id INTEGER NOT NULL);`
	CreateSession         = `INSERT INTO session (sess_id, user_id) VALUES ($1, $2);`
	FindSession           = `SELECT sess_id, user_id FROM session WHERE sess_id = $1;`
	DeleteSession         = `DELETE FROM session WHERE sess_id = $1;`
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
	res, err := a.db.Exec(CreateUser, user.FirstName, user.LastName, user.Email, user.Password)
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
	row := a.db.QueryRow(GetUserByEmail, email)

	var user models.User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres get user: %w", myErr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("postgres get user: %w", err)
	}

	return &user, nil
}

func (a *Adapter) CreateNewUserTable() error {
	_, err := a.db.Exec(CreateUserTable)
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
		return fmt.Errorf("postgres create session: %w", myErr.ErrSessionAlreadyExists)
	}

	return nil
}

func (a *Adapter) FindSession(sessID string) (*models.Session, error) {
	res := a.db.QueryRow(FindSession, sessID)
	var sess models.Session
	err := res.Scan(&sess.ID, &sess.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres find session: %w", myErr.ErrSessionNotFound)
		}
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

	retrying := 10
	i := 1
	log.Printf("try ping:%v", i)
	for err = db.Ping(); err != nil; err = db.Ping() {
		if i >= retrying {
			return nil, fmt.Errorf("postgres connect: %w", err)
		}
		i++
		time.Sleep(1 * time.Second)
		log.Printf("try ping postgresql: %v", i)
	}

	return db, nil
}
