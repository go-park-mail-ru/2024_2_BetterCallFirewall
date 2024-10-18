package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var (
	errMockDB      = errors.New("mock error from db")
	errMockProfile = errors.New("mock profile error")
)

var Posts = []*models.Post{
	{ID: 1, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "post from author 1"}},
	{ID: 2, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "post from author 2"}},
	{ID: 3, Header: models.Header{AuthorID: 3}, PostContent: models.Content{Text: "post from author 3"}},
	{ID: 4, Header: models.Header{AuthorID: 4}, PostContent: models.Content{Text: "post from author 4"}},
	{ID: 5, Header: models.Header{AuthorID: 5}, PostContent: models.Content{Text: "post from author 5"}},
}

type mockDB struct {
	counter uint32
}

func (m *mockDB) Create(ctx context.Context, post *models.Post) (uint32, error) {
	if post.Header.AuthorID == 0 {
		return 0, errMockDB
	}

	return m.counter, nil
}

func (m *mockDB) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	if postID == 0 {
		return nil, errMockDB
	}

	if postID == 1 {
		return &models.Post{ID: postID, Header: models.Header{AuthorID: 0}, PostContent: models.Content{Text: "post with wrong author"}}, nil
	}

	return &models.Post{ID: postID, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "post with real author"}}, nil
}

func (m *mockDB) Update(ctx context.Context, post *models.Post) error {
	if post.ID == 0 {
		return errMockDB
	}

	return nil
}

func (m *mockDB) Delete(ctx context.Context, postID uint32) error {
	if postID == 0 {
		return errMockDB
	}

	return nil
}

func (m *mockDB) GetPosts(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	if lastID > 5 || lastID == 0 {
		return nil, errMockDB
	}
	if lastID == 1 {
		return []*models.Post{{Header: models.Header{AuthorID: 0}}}, nil
	}

	return Posts[:lastID], nil
}

func (m *mockDB) GetFriendsPosts(ctx context.Context, friendsID []uint32, lastID uint32) ([]*models.Post, error) {
	if lastID > 5 || lastID == 0 {
		return nil, errMockDB
	}

	if lastID == 1 {
		return []*models.Post{{Header: models.Header{AuthorID: 0}}}, nil
	}

	var res []*models.Post

	for _, id := range friendsID {
		for i := 0; i <= int(lastID)-1; i++ {
			if Posts[i].Header.AuthorID == id {
				res = append(res, Posts[i])
			}
		}
	}

	return res, nil
}

func (m *mockDB) GetPostAuthor(ctx context.Context, postID uint32) (uint32, error) {
	if postID == 0 {
		return 0, errMockDB
	}

	return postID, nil
}

type profileRepositoryMock struct{}

func (p *profileRepositoryMock) GetHeader(userID uint32) (models.Header, error) {
	if userID == 0 {
		return models.Header{}, errMockProfile
	}

	return models.Header{Author: "Alexey Zemliakov", AuthorID: 1}, nil
}

func (p *profileRepositoryMock) GetFriendsID(userID uint32) ([]uint32, error) {
	if userID == 0 {
		return nil, errMockProfile
	}

	var (
		i   uint32
		ids []uint32
	)
	for i = 1; i <= userID; i++ {
		ids = append(ids, i)
	}

	return ids, nil
}

type TestCaseCreate struct {
	post    *models.Post
	wantID  uint32
	wantErr error
}

var (
	baseId uint32 = 1
	db            = &mockDB{counter: baseId}
	pr            = &profileRepositoryMock{}
	ctx           = context.Background()
)

func TestPostServiceCreate(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
		return
	}

	tests := []TestCaseCreate{
		{post: &models.Post{Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "text of content"}}, wantID: baseId, wantErr: nil},
		{post: &models.Post{Header: models.Header{AuthorID: 0}, PostContent: models.Content{Text: ""}}, wantID: 0, wantErr: errMockDB},
		{post: &models.Post{Header: models.Header{AuthorID: 20}, PostContent: models.Content{Text: "post about VK"}}, wantID: baseId, wantErr: nil},
	}

	for i, tt := range tests {
		id, err := service.Create(ctx, tt.post)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		if id != tt.wantID {
			t.Errorf("#%d: id mismatch: exp=%d got=%d", i, tt.wantID, id)
		}
	}
}

type TestCaseGet struct {
	ID       uint32
	wantPost *models.Post
	wantErr  error
}

func TestPostServiceGet(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseGet{
		{ID: 0, wantPost: nil, wantErr: errMockDB},
		{ID: 1, wantPost: nil, wantErr: errMockProfile},
		{ID: 2, wantPost: &models.Post{ID: 2, PostContent: models.Content{Text: "post with real author"}, Header: models.Header{Author: "Alexey Zemliakov", AuthorID: 1}}, wantErr: nil},
		{ID: 6, wantPost: &models.Post{ID: 6, PostContent: models.Content{Text: "post with real author"}, Header: models.Header{Author: "Alexey Zemliakov", AuthorID: 1}}, wantErr: nil},
	}

	for i, tt := range tests {
		gotPost, err := service.Get(ctx, tt.ID)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		assert.Equal(t, tt.wantPost, gotPost, "#%d: post mismatch:\n exp=%v\n got=%v", i, tt.wantPost, gotPost)
	}
}

type TestCaseDelete struct {
	ID      uint32
	wantErr error
}

func TestPostServiceDelete(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseDelete{
		{ID: 0, wantErr: errMockDB},
		{ID: 1, wantErr: nil},
		{ID: 2, wantErr: nil},
	}

	for i, tt := range tests {
		err := service.Delete(ctx, tt.ID)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
	}
}

type TestCaseUpdate struct {
	post    *models.Post
	wantErr error
}

func TestPostServiceUpdate(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseUpdate{
		{post: &models.Post{ID: 0}, wantErr: errMockDB},
		{post: &models.Post{ID: 1}, wantErr: nil},
		{post: &models.Post{ID: 2}, wantErr: nil},
	}

	for i, tt := range tests {
		err := service.Update(ctx, tt.post)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
	}
}

type TestCaseGetBatch struct {
	lastId    uint32
	wantPosts []*models.Post
	wantErr   error
}

func TestPostServiceGetBatch(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseGetBatch{
		{lastId: 0, wantPosts: nil, wantErr: errMockDB},
		{lastId: 1, wantPosts: nil, wantErr: errMockProfile},
		{lastId: 10, wantPosts: nil, wantErr: errMockDB},
		{lastId: 2, wantPosts: Posts[:2], wantErr: nil},
		{lastId: 5, wantPosts: Posts[:], wantErr: nil},
		{lastId: 4, wantPosts: Posts[:4], wantErr: nil},
	}

	for i, tt := range tests {
		gotPosts, err := service.GetBatch(ctx, tt.lastId)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		assert.Equalf(t, gotPosts, tt.wantPosts, "#%d: post mismatch:\n exp=%v got=%v", i, tt.wantPosts, gotPosts)
	}
}

type TestCaseGetBatchFromFriend struct {
	lastId    uint32
	userId    uint32
	wantPosts []*models.Post
	wantErr   error
}

func TestPostServiceGetBatchFromFriend(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}
	tests := []TestCaseGetBatchFromFriend{
		{lastId: 0, userId: 10, wantPosts: nil, wantErr: errMockDB},
		{lastId: 10, userId: 10, wantPosts: nil, wantErr: errMockDB},
		{lastId: 1, userId: 0, wantPosts: nil, wantErr: errMockProfile},
		{lastId: 1, userId: 10, wantPosts: nil, wantErr: errMockProfile},
		{lastId: 3, userId: 5, wantPosts: Posts[:3], wantErr: nil},
		{lastId: 5, userId: 5, wantPosts: Posts[:], wantErr: nil},
		{lastId: 2, userId: 3, wantPosts: Posts[:2], wantErr: nil},
	}

	for i, tt := range tests {
		posts, err := service.GetBatchFromFriend(ctx, tt.userId, tt.lastId)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		assert.Equalf(t, posts, tt.wantPosts, "#%d: post mismatch:\n exp=%v got=%v", i, tt.wantPosts, posts)
	}
}

type TestCaseGetAuthor struct {
	postID  uint32
	wantID  uint32
	wantErr error
}

func TestPostServiceGetAuthor(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseGetAuthor{
		{postID: 0, wantID: 0, wantErr: errMockDB},
		{postID: 1, wantID: 1, wantErr: nil},
		{postID: 10, wantID: 10, wantErr: nil},
	}

	for i, tt := range tests {
		id, err := service.GetPostAuthorID(ctx, tt.postID)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		if id != tt.wantID {
			t.Errorf("#%d: id mismatch:\n exp=%v got=%v", i, tt.wantID, id)
		}
	}
}
