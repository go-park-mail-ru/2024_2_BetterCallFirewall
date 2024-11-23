package post

import (
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mockSessionManager struct{}

func (m mockSessionManager) Check(s string) (*models.Session, error) {
	return nil, nil
}

func (m mockSessionManager) Create(userID uint32) (*models.Session, error) {
	return nil, nil
}

func (m mockSessionManager) Destroy(sess *models.Session) error {
	return nil
}

type mockPostController struct{}

func (m mockPostController) SetLikeOnPost(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (m mockPostController) DeleteLikeFromPost(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (m mockPostController) GetLikesOnPost(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (m mockPostController) Create(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) GetOne(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) Update(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) Delete(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) GetBatchPosts(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	r := NewRouter(mockPostController{}, mockSessionManager{}, logrus.New())
	assert.NotNil(t, r)
}
