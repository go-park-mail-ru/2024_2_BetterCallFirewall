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
