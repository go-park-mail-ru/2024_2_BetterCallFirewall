package csat

import (
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/csat_api"
)

func NewRequest(id uint32) *csat_api.Request {
	return &csat_api.Request{
		UserID: id,
	}
}

func NewReqWithTime(id uint32, dur time.Duration) *csat_api.RequestWithTime {
	return &csat_api.RequestWithTime{
		UserID:    id,
		SpentTime: int64(dur),
	}
}
