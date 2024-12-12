package stickers

import (
	"net/http"
)

type Controller interface {
	AddNewSticker(w http.ResponseWriter, r *http.Request)
	GetAllStickers(w http.ResponseWriter, r *http.Request)
	GetMineStickers(w http.ResponseWriter, r *http.Request)
}
