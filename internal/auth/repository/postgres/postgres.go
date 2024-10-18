package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type Adapter struct {
	db *sql.DB
}

func NewAdapter(db *sql.DB) *Adapter {
	adapter := &Adapter{
		db: db,
	}
	go adapter.startSessionGC()
	return adapter
}

func (a *Adapter) Create(user *models.User, ctx context.Context) (uint32, error) {
	var id uint32
	err := a.db.QueryRowContext(ctx, CreateUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres create user rows affected: %w", err)
		}
		return 0, fmt.Errorf("postgres create user: %w", err)
	}

	return id, nil
}

func (a *Adapter) GetByEmail(email string, ctx context.Context) (*models.User, error) {
	user := &models.User{}
	err := a.db.QueryRowContext(ctx, GetUserByEmail, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres get user: %w", myErr.ErrUserNotFound)
		}
		return nil, fmt.Errorf("postgres get user: %w", err)
	}

	return user, nil
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
		return fmt.Errorf("postgres create session table: %w", err)
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
		return fmt.Errorf("postgres delete session table: %w", err)
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
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	db.SetMaxOpenConns(10)

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
