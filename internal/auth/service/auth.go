package service

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
)

type UserRepo interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type AuthServiceImpl struct {
	db UserRepo
}

func NewAuthServiceImpl(db UserRepo) *AuthServiceImpl {
	return &AuthServiceImpl{
		db: db,
	}
}

func (a *AuthServiceImpl) Register(user models.User) error {
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

func (a *AuthServiceImpl) validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[\w-.]+@([\w-]+\.)\w{2,4}$`)
	return emailRegex.MatchString(email)
}
