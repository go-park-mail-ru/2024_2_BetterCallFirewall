// Code generated by MockGen. DO NOT EDIT.
// Source: grpc_server.go

// Package post_api is a generated GoMock package.
package post_api

import (
	context "context"
	reflect "reflect"

	models "github.com/2024_2_BetterCallFirewall/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockPostService is a mock of PostService interface.
type MockPostService struct {
	ctrl     *gomock.Controller
	recorder *MockPostServiceMockRecorder
}

// MockPostServiceMockRecorder is the mock recorder for MockPostService.
type MockPostServiceMockRecorder struct {
	mock *MockPostService
}

// NewMockPostService creates a new mock instance.
func NewMockPostService(ctrl *gomock.Controller) *MockPostService {
	mock := &MockPostService{ctrl: ctrl}
	mock.recorder = &MockPostServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostService) EXPECT() *MockPostServiceMockRecorder {
	return m.recorder
}

// GetAuthorsPosts mocks base method.
func (m *MockPostService) GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthorsPosts", ctx, header)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthorsPosts indicates an expected call of GetAuthorsPosts.
func (mr *MockPostServiceMockRecorder) GetAuthorsPosts(ctx, header interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthorsPosts", reflect.TypeOf((*MockPostService)(nil).GetAuthorsPosts), ctx, header)
}
