package file

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

type mockFileController struct{}

func (m mockFileController) Upload(w http.ResponseWriter, r *http.Request) {}

func (m mockFileController) Download(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	r := NewRouter(mockFileController{}, mockSessionManager{}, logrus.New(), &metrics.FileMetrics{})
	assert.NotNil(t, r)
}
