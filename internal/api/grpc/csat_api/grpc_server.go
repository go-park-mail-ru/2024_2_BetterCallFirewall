package csat_api

import (
	"context"
	"time"
)

type CSATService interface {
	NewLike(id uint32)
	NewFriend(id uint32)
	NewMessage(id uint32)
	TimeSpent(id uint32, dur time.Duration)
}

type Adapter struct {
	UnimplementedCsatServiceServer
	service CSATService
}

func NewAdapter(service CSATService) *Adapter {
	return &Adapter{service: service}
}

func (a *Adapter) NewLike(ctx context.Context, req *Request) (*EmptyResponse, error) {
	a.service.NewLike(req.UserID)
	return &EmptyResponse{}, nil
}

func (a *Adapter) NewFriend(ctx context.Context, req *Request) (*EmptyResponse, error) {
	a.service.NewFriend(req.UserID)
	return &EmptyResponse{}, nil
}

func (a *Adapter) NewMessage(ctx context.Context, req *Request) (*EmptyResponse, error) {
	a.service.NewMessage(req.UserID)
	return &EmptyResponse{}, nil
}

func (a *Adapter) TimeSpent(ctx context.Context, req *RequestWithTime) (*EmptyResponse, error) {
	a.service.TimeSpent(req.UserID, time.Duration(req.SpentTime))
	return &EmptyResponse{}, nil
}
