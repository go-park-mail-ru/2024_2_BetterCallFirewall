package csat_api

import (
	"context"
)

type CSATService interface {
	NewLike(id uint32)
	NewFriend(id uint32)
	NewMessage(id uint32)
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
