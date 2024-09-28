package repository

import (
	"time"

	"github.com/brianvoe/gofakeit"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type Repository struct {
	storage  map[string]*models.Post
	havePost bool
}

func NewRepository() *Repository {
	return &Repository{
		storage:  make(map[string]*models.Post),
		havePost: false,
	}
}

func (r *Repository) FakeData(count int) {
	for i := 0; i < count; i++ {
		title := gofakeit.FirstName() + " " + gofakeit.LastName()
		content := gofakeit.Sentence(10)
		date := gofakeit.Date()
		r.storage[title] = &models.Post{
			Header:    title,
			Body:      content,
			CreatedAt: date.Format(time.DateOnly),
		}
	}

	r.havePost = true
}

func (r *Repository) GetAll() ([]*models.Post, error) {
	if !r.havePost {
		return nil, myErr.ErrPostEnd
	}

	res := make([]*models.Post, 0, len(r.storage))
	for _, post := range r.storage {
		res = append(res, post)
	}
	r.havePost = false

	return res, nil
}
