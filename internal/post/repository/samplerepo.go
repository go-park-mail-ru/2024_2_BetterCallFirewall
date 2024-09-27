package repository

import (
	"github.com/brianvoe/gofakeit"

	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type Repository struct {
	storage map[string]*models.Post
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[string]*models.Post),
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
			CreatedAt: date,
		}
	}
}

func (r *Repository) GetAll() []*models.Post {
	res := make([]*models.Post, 0, len(r.storage))
	for _, post := range r.storage {
		res = append(res, post)
	}
	return res
}
