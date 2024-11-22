package community

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

type mockCommunityController struct{}

func (m mockCommunityController) JoinToCommunity(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) LeaveFromCommunity(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) AddAdmin(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) GetAll(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) GetOne(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) Update(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) Delete(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) Create(w http.ResponseWriter, r *http.Request) {}

func (m mockCommunityController) SearchCommunity(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	r := NewRouter(mockCommunityController{}, mockSessionManager{}, logrus.New())
	assert.NotNil(t, r)
}
