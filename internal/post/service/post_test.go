package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/post/entities"
)

var (
	errMockDB      = errors.New("mock error from db")
	errMockProfile = errors.New("mock profile error")
)

var Posts = []*models.Post{
	{ID: 1, AuthorID: 1, PostContent: models.Content{Text: "post from author 1"}},
	{ID: 2, AuthorID: 2, PostContent: models.Content{Text: "post from author 2"}},
	{ID: 3, AuthorID: 3, PostContent: models.Content{Text: "post from author 3"}},
	{ID: 4, AuthorID: 4, PostContent: models.Content{Text: "post from author 4"}},
	{ID: 5, AuthorID: 5, PostContent: models.Content{Text: "post from author 5"}},
}

type mockDB struct {
	counter uint32
}

func (m *mockDB) Create(post *entities.PostDB) (uint32, error) {
	if post.AuthorID == 0 {
		return 0, errMockDB
	}

	return m.counter, nil
}

func (m *mockDB) Get(postID uint32) (*models.Post, error) {
	if postID == 0 {
		return nil, errMockDB
	}

	if postID == 1 {
		return &models.Post{ID: postID, AuthorID: 0, PostContent: models.Content{Text: "post with wrong author"}}, nil
	}

	return &models.Post{ID: postID, AuthorID: 1, PostContent: models.Content{Text: "post with real author"}}, nil
}

func (m *mockDB) Update(post *entities.PostDB) error {
	if post.ID == 0 {
		return errMockDB
	}

	return nil
}

func (m *mockDB) Delete(postID uint32) error {
	if postID == 0 {
		return errMockDB
	}

	return nil
}

func (m *mockDB) CheckAccess(profileID uint32, postID uint32) (bool, error) {
	if profileID == 0 || postID == 0 {
		return false, errMockDB
	}

	if postID != profileID {
		return false, nil
	}

	return true, nil
}

func (m *mockDB) GetPosts(lastID uint32, newRequest bool) ([]*models.Post, error) {
	if lastID > 5 || lastID == 0 {
		return nil, errMockDB
	}
	if lastID == 1 {
		return []*models.Post{{AuthorID: 0}}, nil
	}

	if newRequest {
		return Posts, nil
	}

	return Posts[:lastID], nil
}

func (m *mockDB) GetFriendsPosts(friendsID []uint32, lastID uint32, newRequest bool) ([]*models.Post, error) {
	if lastID > 5 || lastID == 0 {
		return nil, errMockDB
	}

	if lastID == 1 {
		return []*models.Post{{AuthorID: 0}}, nil
	}

	var res []*models.Post

	if newRequest {
		for _, id := range friendsID {
			for _, post := range Posts {
				if post.AuthorID == id {
					res = append(res, post)
				}
			}
		}
		return res, nil
	}

	for _, id := range friendsID {
		for i := int(lastID) - 1; i >= 0; i-- {
			if Posts[i].AuthorID == id {
				res = append(res, Posts[i])
			}
		}
	}

	return res, nil
}

type profileRepositoryMock struct{}

func (p *profileRepositoryMock) GetHeader(userID uint32) (models.Header, error) {
	if userID == 0 {
		return models.Header{}, errMockProfile
	}

	return models.Header{Author: "Alexey Zemliakov"}, nil
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
)

func TestPostServiceCreate(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
		return
	}

	tests := []TestCaseCreate{
		{post: &models.Post{AuthorID: 1, PostContent: models.Content{Text: "text of content"}}, wantID: baseId, wantErr: nil},
		{post: &models.Post{AuthorID: 0, PostContent: models.Content{Text: ""}}, wantID: 0, wantErr: errMockDB},
		{post: &models.Post{AuthorID: 20, PostContent: models.Content{Text: "post about VK"}}, wantID: baseId, wantErr: nil},
	}

	for i, tt := range tests {
		id, err := service.Create(tt.post)
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
		{ID: 2, wantPost: &models.Post{ID: 2, AuthorID: 1, PostContent: models.Content{Text: "post with real author"}, Header: models.Header{Author: "Alexey Zemliakov"}}, wantErr: nil},
		{ID: 6, wantPost: &models.Post{ID: 6, AuthorID: 1, PostContent: models.Content{Text: "post with real author"}, Header: models.Header{Author: "Alexey Zemliakov"}}, wantErr: nil},
	}

	for i, tt := range tests {
		gotPost, err := service.Get(tt.ID)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		if !reflect.DeepEqual(gotPost, tt.wantPost) {
			t.Errorf("#%d: post mismatch:\n exp=%v\n got=%v", i, tt.wantPost, gotPost)
		}
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
		err := service.Delete(tt.ID)
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
		err := service.Update(tt.post)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
	}
}

type TestCaseCheckAccess struct {
	postID  uint32
	userID  uint32
	want    bool
	wantErr error
}

func TestPostServiceCheckUserAccess(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseCheckAccess{
		{postID: 0, userID: 0, wantErr: errMockDB, want: false},
		{postID: 0, userID: 10, wantErr: errMockDB, want: false},
		{postID: 10, userID: 0, wantErr: errMockDB, want: false},
		{postID: 10, userID: 1, wantErr: nil, want: false},
		{postID: 10, userID: 10, wantErr: nil, want: true},
	}

	for i, tt := range tests {
		ok, err := service.CheckUserAccess(tt.postID, tt.userID)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		if ok != tt.want {
			t.Errorf("#%d: ok mismatch: exp=%v got=%v", i, tt.want, ok)
		}
	}
}

type TestCaseGetBatch struct {
	lastId     uint32
	newRequest bool
	wantPosts  []*models.Post
	wantErr    error
}

func TestPostServiceGetBatch(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}

	tests := []TestCaseGetBatch{
		{lastId: 0, newRequest: false, wantPosts: nil, wantErr: errMockDB},
		{lastId: 1, newRequest: true, wantPosts: nil, wantErr: errMockProfile},
		{lastId: 10, newRequest: true, wantPosts: nil, wantErr: errMockDB},
		{lastId: 2, newRequest: false, wantPosts: Posts[:2], wantErr: nil},
		{lastId: 2, newRequest: true, wantPosts: Posts[:], wantErr: nil},
		{lastId: 4, newRequest: false, wantPosts: Posts[:4], wantErr: nil},
	}

	for i, tt := range tests {
		gotPosts, err := service.GetBatch(tt.lastId, tt.newRequest)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		if !reflect.DeepEqual(gotPosts, tt.wantPosts) {
			t.Errorf("#%d: post mismatch:\n exp=%v got=%v", i, tt.wantPosts, gotPosts)
		}
	}
}

type TestCaseGetBatchFromFriend struct {
	lastId     uint32
	userId     uint32
	newRequest bool
	wantPosts  []*models.Post
	wantErr    error
}

func TestPostServiceGetBatchFromFriend(t *testing.T) {
	service := NewPostServiceImpl(db, pr)
	if service == nil {
		t.Fatal("service is nil")
	}
	tests := []TestCaseGetBatchFromFriend{
		{lastId: 0, userId: 10, newRequest: false, wantPosts: nil, wantErr: errMockDB},
		{lastId: 10, userId: 10, newRequest: false, wantPosts: nil, wantErr: errMockDB},
		{lastId: 1, userId: 0, newRequest: true, wantPosts: nil, wantErr: errMockProfile},
		{lastId: 1, userId: 10, newRequest: true, wantPosts: nil, wantErr: errMockProfile},
		{lastId: 3, userId: 5, newRequest: false, wantPosts: Posts[:3], wantErr: nil},
		{lastId: 2, userId: 5, newRequest: true, wantPosts: Posts[:], wantErr: nil},
		{lastId: 2, userId: 3, newRequest: true, wantPosts: Posts[:3], wantErr: nil},
	}

	for i, tt := range tests {
		posts, err := service.GetBatchFromFriend(tt.userId, tt.lastId, tt.newRequest)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("#%d: error mismatch: exp=%v got=%v", i, tt.wantErr, err)
		}
		if !reflect.DeepEqual(posts, tt.wantPosts) {
			t.Errorf("#%d: post mismatch:\n exp=%v got=%v", i, tt.wantPosts, posts)
		}
	}
}
