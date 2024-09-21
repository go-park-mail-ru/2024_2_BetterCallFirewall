package repository

import (
	"errors"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type SampleDBImpl struct {
	repo    map[string]*User
	counter int
}

func NewSampleDB() *SampleDBImpl {
	return &SampleDBImpl{
		repo:    make(map[string]*User),
		counter: 0,
	}
}

func (s *SampleDBImpl) Create(user models.User) error {
	u := FromServiceUser(user, s.counter)
	s.counter++
	s.repo[u.Email] = &u
	return nil
}

func (s *SampleDBImpl) GetByEmail(email string) (models.User, error) {
	u, ok := s.repo[email]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return u.ToServiceUser(), nil
}
