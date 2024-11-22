// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	models "github.com/2024_2_BetterCallFirewall/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// CheckAccess mocks base method.
func (m *MockRepo) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccess", ctx, communityID, userID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckAccess indicates an expected call of CheckAccess.
func (mr *MockRepoMockRecorder) CheckAccess(ctx, communityID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccess", reflect.TypeOf((*MockRepo)(nil).CheckAccess), ctx, communityID, userID)
}

// Create mocks base method.
func (m *MockRepo) Create(ctx context.Context, community *models.Community, author uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, community, author)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepoMockRecorder) Create(ctx, community, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), ctx, community, author)
}

// Delete mocks base method.
func (m *MockRepo) Delete(ctx context.Context, id uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepo)(nil).Delete), ctx, id)
}

// GetBatch mocks base method.
func (m *MockRepo) GetBatch(ctx context.Context, lastID uint32) ([]*models.CommunityCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatch", ctx, lastID)
	ret0, _ := ret[0].([]*models.CommunityCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatch indicates an expected call of GetBatch.
func (mr *MockRepoMockRecorder) GetBatch(ctx, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatch", reflect.TypeOf((*MockRepo)(nil).GetBatch), ctx, lastID)
}

// GetOne mocks base method.
func (m *MockRepo) GetOne(ctx context.Context, id uint32) (*models.Community, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", ctx, id)
	ret0, _ := ret[0].(*models.Community)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockRepoMockRecorder) GetOne(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockRepo)(nil).GetOne), ctx, id)
}

// JoinCommunity mocks base method.
func (m *MockRepo) JoinCommunity(ctx context.Context, communityId, author uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinCommunity", ctx, communityId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinCommunity indicates an expected call of JoinCommunity.
func (mr *MockRepoMockRecorder) JoinCommunity(ctx, communityId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinCommunity", reflect.TypeOf((*MockRepo)(nil).JoinCommunity), ctx, communityId, author)
}

// LeaveCommunity mocks base method.
func (m *MockRepo) LeaveCommunity(ctx context.Context, communityId, author uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaveCommunity", ctx, communityId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveCommunity indicates an expected call of LeaveCommunity.
func (mr *MockRepoMockRecorder) LeaveCommunity(ctx, communityId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveCommunity", reflect.TypeOf((*MockRepo)(nil).LeaveCommunity), ctx, communityId, author)
}

// NewAdmin mocks base method.
func (m *MockRepo) NewAdmin(ctx context.Context, communityId, author uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewAdmin", ctx, communityId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewAdmin indicates an expected call of NewAdmin.
func (mr *MockRepoMockRecorder) NewAdmin(ctx, communityId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewAdmin", reflect.TypeOf((*MockRepo)(nil).NewAdmin), ctx, communityId, author)
}

// Search mocks base method.
func (m *MockRepo) Search(ctx context.Context, query string, lastID uint32) ([]*models.CommunityCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query, lastID)
	ret0, _ := ret[0].([]*models.CommunityCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockRepoMockRecorder) Search(ctx, query, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockRepo)(nil).Search), ctx, query, lastID)
}

// Update mocks base method.
func (m *MockRepo) Update(ctx context.Context, community *models.Community) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, community)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepoMockRecorder) Update(ctx, community interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepo)(nil).Update), ctx, community)
}
