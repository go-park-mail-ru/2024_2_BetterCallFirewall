// Code generated by MockGen. DO NOT EDIT.
// Source: controller.go

// Package controller is a generated GoMock package.
package controller

import (
	context "context"
	http "net/http"
	reflect "reflect"

	models "github.com/2024_2_BetterCallFirewall/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// Mockresponder is a mock of responder interface.
type Mockresponder struct {
	ctrl     *gomock.Controller
	recorder *MockresponderMockRecorder
}

// MockresponderMockRecorder is the mock recorder for Mockresponder.
type MockresponderMockRecorder struct {
	mock *Mockresponder
}

// NewMockresponder creates a new mock instance.
func NewMockresponder(ctrl *gomock.Controller) *Mockresponder {
	mock := &Mockresponder{ctrl: ctrl}
	mock.recorder = &MockresponderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockresponder) EXPECT() *MockresponderMockRecorder {
	return m.recorder
}

// ErrorBadRequest mocks base method.
func (m *Mockresponder) ErrorBadRequest(w http.ResponseWriter, err error, requestID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ErrorBadRequest", w, err, requestID)
}

// ErrorBadRequest indicates an expected call of ErrorBadRequest.
func (mr *MockresponderMockRecorder) ErrorBadRequest(w, err, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorBadRequest", reflect.TypeOf((*Mockresponder)(nil).ErrorBadRequest), w, err, requestID)
}

// ErrorInternal mocks base method.
func (m *Mockresponder) ErrorInternal(w http.ResponseWriter, err error, requestID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ErrorInternal", w, err, requestID)
}

// ErrorInternal indicates an expected call of ErrorInternal.
func (mr *MockresponderMockRecorder) ErrorInternal(w, err, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorInternal", reflect.TypeOf((*Mockresponder)(nil).ErrorInternal), w, err, requestID)
}

// LogError mocks base method.
func (m *Mockresponder) LogError(err error, requestID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogError", err, requestID)
}

// LogError indicates an expected call of LogError.
func (mr *MockresponderMockRecorder) LogError(err, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogError", reflect.TypeOf((*Mockresponder)(nil).LogError), err, requestID)
}

// OutputJSON mocks base method.
func (m *Mockresponder) OutputJSON(w http.ResponseWriter, data any, requestID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OutputJSON", w, data, requestID)
}

// OutputJSON indicates an expected call of OutputJSON.
func (mr *MockresponderMockRecorder) OutputJSON(w, data, requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OutputJSON", reflect.TypeOf((*Mockresponder)(nil).OutputJSON), w, data, requestID)
}

// OutputNoMoreContentJSON mocks base method.
func (m *Mockresponder) OutputNoMoreContentJSON(w http.ResponseWriter, requestId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OutputNoMoreContentJSON", w, requestId)
}

// OutputNoMoreContentJSON indicates an expected call of OutputNoMoreContentJSON.
func (mr *MockresponderMockRecorder) OutputNoMoreContentJSON(w, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OutputNoMoreContentJSON", reflect.TypeOf((*Mockresponder)(nil).OutputNoMoreContentJSON), w, requestId)
}

// MockcommunityService is a mock of communityService interface.
type MockcommunityService struct {
	ctrl     *gomock.Controller
	recorder *MockcommunityServiceMockRecorder
}

// MockcommunityServiceMockRecorder is the mock recorder for MockcommunityService.
type MockcommunityServiceMockRecorder struct {
	mock *MockcommunityService
}

// NewMockcommunityService creates a new mock instance.
func NewMockcommunityService(ctrl *gomock.Controller) *MockcommunityService {
	mock := &MockcommunityService{ctrl: ctrl}
	mock.recorder = &MockcommunityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcommunityService) EXPECT() *MockcommunityServiceMockRecorder {
	return m.recorder
}

// AddAdmin mocks base method.
func (m *MockcommunityService) AddAdmin(ctx context.Context, communityId, author uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAdmin", ctx, communityId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAdmin indicates an expected call of AddAdmin.
func (mr *MockcommunityServiceMockRecorder) AddAdmin(ctx, communityId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAdmin", reflect.TypeOf((*MockcommunityService)(nil).AddAdmin), ctx, communityId, author)
}

// CheckAccess mocks base method.
func (m *MockcommunityService) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccess", ctx, communityID, userID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckAccess indicates an expected call of CheckAccess.
func (mr *MockcommunityServiceMockRecorder) CheckAccess(ctx, communityID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccess", reflect.TypeOf((*MockcommunityService)(nil).CheckAccess), ctx, communityID, userID)
}

// Create mocks base method.
func (m *MockcommunityService) Create(ctx context.Context, community *models.Community, authorID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, community, authorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockcommunityServiceMockRecorder) Create(ctx, community, authorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockcommunityService)(nil).Create), ctx, community, authorID)
}

// Delete mocks base method.
func (m *MockcommunityService) Delete(ctx context.Context, id uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockcommunityServiceMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockcommunityService)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockcommunityService) Get(ctx context.Context, userID, lastID uint32) ([]*models.CommunityCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID, lastID)
	ret0, _ := ret[0].([]*models.CommunityCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockcommunityServiceMockRecorder) Get(ctx, userID, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockcommunityService)(nil).Get), ctx, userID, lastID)
}

// GetOne mocks base method.
func (m *MockcommunityService) GetOne(ctx context.Context, id, userID uint32) (*models.Community, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", ctx, id, userID)
	ret0, _ := ret[0].(*models.Community)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockcommunityServiceMockRecorder) GetOne(ctx, id, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockcommunityService)(nil).GetOne), ctx, id, userID)
}

// JoinCommunity mocks base method.
func (m *MockcommunityService) JoinCommunity(ctx context.Context, communityId, author uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinCommunity", ctx, communityId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinCommunity indicates an expected call of JoinCommunity.
func (mr *MockcommunityServiceMockRecorder) JoinCommunity(ctx, communityId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinCommunity", reflect.TypeOf((*MockcommunityService)(nil).JoinCommunity), ctx, communityId, author)
}

// LeaveCommunity mocks base method.
func (m *MockcommunityService) LeaveCommunity(ctx context.Context, communityId, author uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaveCommunity", ctx, communityId, author)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveCommunity indicates an expected call of LeaveCommunity.
func (mr *MockcommunityServiceMockRecorder) LeaveCommunity(ctx, communityId, author interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveCommunity", reflect.TypeOf((*MockcommunityService)(nil).LeaveCommunity), ctx, communityId, author)
}

// Search mocks base method.
func (m *MockcommunityService) Search(ctx context.Context, query string, userID, lastID uint32) ([]*models.CommunityCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query, userID, lastID)
	ret0, _ := ret[0].([]*models.CommunityCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockcommunityServiceMockRecorder) Search(ctx, query, userID, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockcommunityService)(nil).Search), ctx, query, userID, lastID)
}

// Update mocks base method.
func (m *MockcommunityService) Update(ctx context.Context, id uint32, community *models.Community) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, community)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockcommunityServiceMockRecorder) Update(ctx, id, community interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockcommunityService)(nil).Update), ctx, id, community)
}
