package repository

import (
	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type SampleDBImpl struct {
	repo    map[string]*models.User
	counter int
}

func NewSampleDB() *SampleDBImpl {
	return &SampleDBImpl{
		repo:    make(map[string]*models.User),
		counter: 0,
	}
}

func (s *SampleDBImpl) Create(user *models.User) error {
	_, err := s.GetByEmail(user.Email)
	if err != myErr.ErrUserNotFound {
		if err != nil {
			return err
		}
		return myErr.ErrUserAlreadyExists
	}
	s.counter++
	user.ID = s.counter
	s.repo[user.Email] = user
	return nil
}

func (s *SampleDBImpl) GetByEmail(email string) (*models.User, error) {
	u, ok := s.repo[email]
	if !ok {
		return nil, myErr.ErrUserNotFound
	}
	return u, nil
}
