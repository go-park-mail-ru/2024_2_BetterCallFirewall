// Code generated by MockGen. DO NOT EDIT.
// Source: grpc_server.go

// Package community_api is a generated GoMock package.
package community_api

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommunityService is a mock of CommunityService interface.
type MockCommunityService struct {
	ctrl     *gomock.Controller
	recorder *MockCommunityServiceMockRecorder
}

// MockCommunityServiceMockRecorder is the mock recorder for MockCommunityService.
type MockCommunityServiceMockRecorder struct {
	mock *MockCommunityService
}

// NewMockCommunityService creates a new mock instance.
func NewMockCommunityService(ctrl *gomock.Controller) *MockCommunityService {
	mock := &MockCommunityService{ctrl: ctrl}
	mock.recorder = &MockCommunityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunityService) EXPECT() *MockCommunityServiceMockRecorder {
	return m.recorder
}

// CheckAccess mocks base method.
func (m *MockCommunityService) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccess", ctx, communityID, userID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckAccess indicates an expected call of CheckAccess.
func (mr *MockCommunityServiceMockRecorder) CheckAccess(ctx, communityID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccess", reflect.TypeOf((*MockCommunityService)(nil).CheckAccess), ctx, communityID, userID)
}
