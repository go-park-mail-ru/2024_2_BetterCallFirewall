package stickers

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

type mockStickerController struct{}

func (m mockStickerController) GetAllStickers(w http.ResponseWriter, r *http.Request) {}

func (m mockStickerController) AddNewSticker(w http.ResponseWriter, r *http.Request) {}

func (m mockStickerController) GetMineStickers(w http.ResponseWriter, r *http.Request) {}

func TestNewRouter(t *testing.T) {
	r := NewRouter(mockStickerController{}, mockSessionManager{}, logrus.New())
	assert.NotNil(t, r)
}
