package chat

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

type mockChatController struct{}

func (m mockChatController) SetConnection(w http.ResponseWriter, r *http.Request) {}

func (m mockChatController) GetAllChats(w http.ResponseWriter, r *http.Request) {}

func (m mockChatController) GetChat(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	r := NewRouter(mockChatController{}, mockSessionManager{}, logrus.New(), &metrics.HttpMetrics{})
	assert.NotNil(t, r)
}
