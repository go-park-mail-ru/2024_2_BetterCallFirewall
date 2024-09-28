package service

import (
	"reflect"
	"sort"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type mockDB struct {
	container map[string]*models.Post
	hasData   bool
}

func newMockDB() *mockDB {
	return &mockDB{
		container: make(map[string]*models.Post),
		hasData:   false,
	}
}

func (m *mockDB) Save(post *models.Post) {
	m.container[post.Header] = post
	m.hasData = true
}

func (m *mockDB) clearAll() {
	clear(m.container)
	m.hasData = false
}

func (m *mockDB) GetAll() ([]*models.Post, error) {
	if !m.hasData {
		return nil, myErr.ErrPostEnd
	}
	var posts []*models.Post
	for _, post := range m.container {
		posts = append(posts, post)
	}
	m.clearAll()
	return posts, nil
}

type TestCase struct {
	posts   []*models.Post
	want    []*models.Post
	wantErr error
}

func TestPostService(t *testing.T) {
	db := newMockDB()
	service := NewPostServiceImpl(db)

	tests := []TestCase{
		{posts: nil, want: nil, wantErr: myErr.ErrPostEnd},
		{
			posts: []*models.Post{{Header: "Header", Body: "Body", CreatedAt: "2012-10-30"}},
			want:  []*models.Post{{Header: "Header", Body: "Body", CreatedAt: "2012-10-30"}},
		},
		{
			posts: []*models.Post{
				{Header: "Some header", Body: "good text", CreatedAt: "2012-10-30"},
				{Header: "A header", Body: "text about fish", CreatedAt: "2012-10-30"},
			},
			want: []*models.Post{
				{Header: "A header", Body: "text about fish", CreatedAt: "2012-10-30"},
				{Header: "Some header", Body: "good text", CreatedAt: "2012-10-30"},
			},
		},
	}

	for _, test := range tests {
		for _, post := range test.posts {
			db.Save(post)
		}
		got, _ := service.GetAll()
		sort.Slice(got, func(i, j int) bool {
			return got[i].Header < got[j].Header
		})
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetAll() = %v, want %v", got, test.want)
		}
	}
}
