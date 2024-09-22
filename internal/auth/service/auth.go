package service

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type DB interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type AuthServiceImpl struct {
	db DB
}

func NewAuthServiceImpl(db DB) *AuthServiceImpl {
	return &AuthServiceImpl{
		db: db,
	}
}

func (a *AuthServiceImpl) Register(user models.User) error {
	_, err := a.db.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, myErr.ErrUserNotFound) {
		return fmt.Errorf("registration: %w", err)
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("registration: %w", err)
	}
	user.Password = string(hashPassword)
	err = a.db.Create(&user)
	if err != nil {
		return fmt.Errorf("registration: %w", err)
	}

	return nil
}
