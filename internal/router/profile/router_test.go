package profile

import (
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
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

type mockProfileController struct{}

func (m mockProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetProfileById(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAll(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) UpdateProfile(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) DeleteProfile(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetHeader(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) SendFriendReq(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) AcceptFriendReq(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) Unsubscribe(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) RemoveFromFriends(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAllFriends(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAllSubs(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {}

func (m mockProfileController) GetCommunitySubs(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	r := NewRouter(mockProfileController{}, mockSessionManager{}, logrus.New(), &metrics.HttpMetrics{})
	assert.NotNil(t, r)
}
