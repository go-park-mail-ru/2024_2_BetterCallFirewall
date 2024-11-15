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

// CheckAccessToCommunity mocks base method.
func (m *MockPostService) CheckAccessToCommunity(ctx context.Context, userID, communityID uint32) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAccessToCommunity", ctx, userID, communityID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckAccessToCommunity indicates an expected call of CheckAccessToCommunity.
func (mr *MockPostServiceMockRecorder) CheckAccessToCommunity(ctx, userID, communityID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAccessToCommunity", reflect.TypeOf((*MockPostService)(nil).CheckAccessToCommunity), ctx, userID, communityID)
}

// Create mocks base method.
func (m *MockPostService) Create(ctx context.Context, post *models.Post) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, post)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPostServiceMockRecorder) Create(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPostService)(nil).Create), ctx, post)
}

// CreateCommunityPost mocks base method.
func (m *MockPostService) CreateCommunityPost(ctx context.Context, post *models.Post) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommunityPost", ctx, post)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCommunityPost indicates an expected call of CreateCommunityPost.
func (mr *MockPostServiceMockRecorder) CreateCommunityPost(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommunityPost", reflect.TypeOf((*MockPostService)(nil).CreateCommunityPost), ctx, post)
}

// Delete mocks base method.
func (m *MockPostService) Delete(ctx context.Context, postID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPostServiceMockRecorder) Delete(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPostService)(nil).Delete), ctx, postID)
}

// Get mocks base method.
func (m *MockPostService) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, postID)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPostServiceMockRecorder) Get(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPostService)(nil).Get), ctx, postID)
}

// GetBatch mocks base method.
func (m *MockPostService) GetBatch(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatch", ctx, lastID)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatch indicates an expected call of GetBatch.
func (mr *MockPostServiceMockRecorder) GetBatch(ctx, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatch", reflect.TypeOf((*MockPostService)(nil).GetBatch), ctx, lastID)
}

// GetBatchFromFriend mocks base method.
func (m *MockPostService) GetBatchFromFriend(ctx context.Context, userID, lastID uint32) ([]*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchFromFriend", ctx, userID, lastID)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchFromFriend indicates an expected call of GetBatchFromFriend.
func (mr *MockPostServiceMockRecorder) GetBatchFromFriend(ctx, userID, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchFromFriend", reflect.TypeOf((*MockPostService)(nil).GetBatchFromFriend), ctx, userID, lastID)
}

// GetCommunityPost mocks base method.
func (m *MockPostService) GetCommunityPost(ctx context.Context, communityID, lastID uint32) ([]*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommunityPost", ctx, communityID, lastID)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommunityPost indicates an expected call of GetCommunityPost.
func (mr *MockPostServiceMockRecorder) GetCommunityPost(ctx, communityID, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommunityPost", reflect.TypeOf((*MockPostService)(nil).GetCommunityPost), ctx, communityID, lastID)
}

// GetPostAuthorID mocks base method.
func (m *MockPostService) GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostAuthorID", ctx, postID)
	ret0, _ := ret[0].(uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostAuthorID indicates an expected call of GetPostAuthorID.
func (mr *MockPostServiceMockRecorder) GetPostAuthorID(ctx, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostAuthorID", reflect.TypeOf((*MockPostService)(nil).GetPostAuthorID), ctx, postID)
}

// Update mocks base method.
func (m *MockPostService) Update(ctx context.Context, post *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, post)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPostServiceMockRecorder) Update(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPostService)(nil).Update), ctx, post)
}

// MockResponder is a mock of Responder interface.
type MockResponder struct {
	ctrl     *gomock.Controller
	recorder *MockResponderMockRecorder
}

// MockResponderMockRecorder is the mock recorder for MockResponder.
type MockResponderMockRecorder struct {
	mock *MockResponder
}

// NewMockResponder creates a new mock instance.
func NewMockResponder(ctrl *gomock.Controller) *MockResponder {
	mock := &MockResponder{ctrl: ctrl}
	mock.recorder = &MockResponderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResponder) EXPECT() *MockResponderMockRecorder {
	return m.recorder
}

// ErrorBadRequest mocks base method.
func (m *MockResponder) ErrorBadRequest(w http.ResponseWriter, err error, requestId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ErrorBadRequest", w, err, requestId)
}

// ErrorBadRequest indicates an expected call of ErrorBadRequest.
func (mr *MockResponderMockRecorder) ErrorBadRequest(w, err, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorBadRequest", reflect.TypeOf((*MockResponder)(nil).ErrorBadRequest), w, err, requestId)
}

// ErrorInternal mocks base method.
func (m *MockResponder) ErrorInternal(w http.ResponseWriter, err error, requestId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ErrorInternal", w, err, requestId)
}

// ErrorInternal indicates an expected call of ErrorInternal.
func (mr *MockResponderMockRecorder) ErrorInternal(w, err, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorInternal", reflect.TypeOf((*MockResponder)(nil).ErrorInternal), w, err, requestId)
}

// LogError mocks base method.
func (m *MockResponder) LogError(err error, requestId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogError", err, requestId)
}

// LogError indicates an expected call of LogError.
func (mr *MockResponderMockRecorder) LogError(err, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogError", reflect.TypeOf((*MockResponder)(nil).LogError), err, requestId)
}

// OutputJSON mocks base method.
func (m *MockResponder) OutputJSON(w http.ResponseWriter, data any, requestId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OutputJSON", w, data, requestId)
}

// OutputJSON indicates an expected call of OutputJSON.
func (mr *MockResponderMockRecorder) OutputJSON(w, data, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OutputJSON", reflect.TypeOf((*MockResponder)(nil).OutputJSON), w, data, requestId)
}

// OutputNoMoreContentJSON mocks base method.
func (m *MockResponder) OutputNoMoreContentJSON(w http.ResponseWriter, requestId string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OutputNoMoreContentJSON", w, requestId)
}

// OutputNoMoreContentJSON indicates an expected call of OutputNoMoreContentJSON.
func (mr *MockResponderMockRecorder) OutputNoMoreContentJSON(w, requestId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OutputNoMoreContentJSON", reflect.TypeOf((*MockResponder)(nil).OutputNoMoreContentJSON), w, requestId)
}
