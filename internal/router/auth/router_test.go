package auth

import (
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mockController struct{}

func (m mockController) Register(w http.ResponseWriter, r *http.Request) {}

func (m mockController) Auth(w http.ResponseWriter, r *http.Request) {}

func (m mockController) Logout(w http.ResponseWriter, r *http.Request) {}

type mockMiddleware struct{}

func (m mockMiddleware) Check(str string) (*models.Session, error) { return nil, nil }

func (m mockMiddleware) Create(userID uint32) (*models.Session, error) { return nil, nil }

func (m mockMiddleware) Destroy(sess *models.Session) error { return nil }

func TestNewRouter(t *testing.T) {
	router := NewRouter(mockController{}, mockMiddleware{}, logrus.New())
	assert.NotNil(t, router)
}
