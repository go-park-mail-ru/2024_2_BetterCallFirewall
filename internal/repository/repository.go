package repository

import (
	"github.com/2024_2_BetterCallFirewall/internal/custom_error"
	"github.com/2024_2_BetterCallFirewall/internal/models"
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
	if err != custom_error.ErrUserNotFound {
		if err != nil {
			return err
		}
		return custom_error.ErrUserAlreadyExists
	}
	s.counter++
	user.ID = s.counter
	s.repo[user.Email] = user
	return nil
}

func (s *SampleDBImpl) GetByEmail(email string) (*models.User, error) {
	u, ok := s.repo[email]
	if !ok {
		return nil, custom_error.ErrUserNotFound
	}
	return u, nil
}
