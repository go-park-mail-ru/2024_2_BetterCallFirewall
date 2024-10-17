package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

const (
	CreateUserTable       = `CREATE TABLE IF NOT EXISTS person (id INT PRIMARY KEY, first_name TEXT NOT NULL CONSTRAINT first_name_length CHECK (CHAR_LENGTH(first_name) <= 30), last_name TEXT NOT NULL CONSTRAINT last_name_length CHECK (CHAR_LENGTH(last_name) <= 30), email TEXT NOT NULL UNIQUE NOT NULL CONSTRAINT email_length CHECK (CHAR_LENGTH(email) <= 50), password TEXT NOT NULL CONSTRAINT password_length CHECK (CHAR_LENGTH(password) <= 61));`
	CreateUser            = `INSERT INTO person (id, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (email) DO NOTHING;`
	GetUserByEmail        = `SELECT id, first_name, last_name, email, password FROM person WHERE email = $1;`
	CreateNewSessionTable = `CREATE TABLE IF NOT EXISTS session (id SERIAL PRIMARY KEY, sess_id TEXT NOT NULL, user_id INTEGER NOT NULL UNIQUE, created_at BIGINT NOT NULL);`
	CreateSession         = `INSERT INTO session (sess_id, user_id, created_at) VALUES ($1, $2, $3) ON CONFLICT(user_id) DO UPDATE SET sess_id = EXCLUDED.sess_id, created_at = EXCLUDED.created_at;`
	FindSession           = `SELECT sess_id, user_id, created_at FROM session WHERE sess_id = $1;`
	DeleteSession         = `DELETE FROM session WHERE sess_id = $1;`
	DeleteOutdatedSession = `DELETE FROM session WHERE created_at <= $1;`
)

type Adapter struct {
	db      *sql.DB
	counter uint32
}

func NewAdapter(db *sql.DB) *Adapter {
	adapter := &Adapter{
		db:      db,
		counter: 1,
	}
	go adapter.startSessionGC()
	return adapter
}

func (a *Adapter) Create(user *models.User) (uint32, error) {
	res, err := a.db.Exec(CreateUser, a.counter, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return 0, fmt.Errorf("postgres create user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("postgres create user rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("postgres create user: %w", myErr.ErrUserAlreadyExists)
	}
	a.counter++

	return a.counter - 1, nil
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
	res, err := a.db.Exec(CreateSession, sess.ID, sess.UserID, sess.CreatedAt)
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
	err := res.Scan(&sess.ID, &sess.UserID, &sess.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres find session: %w", myErr.ErrSessionNotFound)
		}
		return nil, fmt.Errorf("postgres find session table: %w", err)
	}

	return &sess, nil
}

func (a *Adapter) DestroySession(sessID string) error {
	res, err := a.db.Exec(DeleteSession, sessID)
	if err != nil {
		return fmt.Errorf("postgres delete session table: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres delete session rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("postgres delete session: %w", myErr.ErrSessionNotFound)
	}

	return nil
}

func (a *Adapter) destroyOutdatedSession() error {
	destroyTime := time.Now().Add(-24 * time.Hour).Unix()
	_, err := a.db.Exec(DeleteOutdatedSession, destroyTime)
	if err != nil {
		return fmt.Errorf("postgres destroy outdated session table: %w", err)
	}

	return nil
}

func (a *Adapter) startSessionGC() {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		<-ticker.C
		err := a.destroyOutdatedSession()
		if err != nil {
			log.Println(err)
		}
	}
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
