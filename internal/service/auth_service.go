package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
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
	u, err := a.db.GetByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("registration: %w", err)
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("registration: %w", err)
	}
	u.Password = string(hashPassword)
	return a.db.Create(&user)
}

func (a *AuthServiceImpl) Auth(user models.User) error {
	//TODO implement this method
	return nil
}

func (a *AuthServiceImpl) VerifyToken(token string) error {
	//TODO implement this method
	return nil
}
