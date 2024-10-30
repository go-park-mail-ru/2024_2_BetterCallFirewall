package router

import (
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type MockAuthController struct{}

func (m MockAuthController) Register(w http.ResponseWriter, r *http.Request) {}

func (m MockAuthController) Auth(w http.ResponseWriter, r *http.Request) {}

func (m MockAuthController) Logout(w http.ResponseWriter, r *http.Request) {}

type mockProfileController struct{}

func (m mockProfileController) GetProfileById(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAll(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) UpdateProfile(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) DeleteProfile(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) SendFriendReq(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) AcceptFriendReq(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) RemoveFromFriends(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAllFriends(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) Unsubscribe(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAllSubs(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {}

type mockPostController struct{}

func (m mockPostController) Create(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) GetOne(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) Update(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) Delete(w http.ResponseWriter, r *http.Request) {}

func (m mockPostController) GetBatchPosts(w http.ResponseWriter, r *http.Request) {}

type mockMiddleware struct{}

func (m mockMiddleware) Check(r *http.Request) (*models.Session, error) {
	return nil, nil
}

func (m mockMiddleware) Create(w http.ResponseWriter, userID uint32) (*models.Session, error) {
	return nil, nil
}

func (m mockMiddleware) Destroy(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func TestNewRouter(t *testing.T) {
	router := NewRouter(MockAuthController{}, mockProfileController{}, mockPostController{}, mockMiddleware{}, logrus.New())
	assert.NotNil(t, router)
}