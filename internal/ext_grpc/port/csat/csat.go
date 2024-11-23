package csat

import (
	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/csat_api"
)

func NewRequest(id uint32) *csat_api.Request {
	return &csat_api.Request{
		UserID: id,
	}
}
