package service

import (
	"reflect"
	"sort"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/post/models"
)

type mockDB struct {
	container map[string]*models.Post
}

func newMockDB() *mockDB {
	return &mockDB{
		container: make(map[string]*models.Post),
	}
}

func (m *mockDB) Save(post *models.Post) {
	m.container[post.Header] = post
}

func (m *mockDB) ClearAll() {
	clear(m.container)
}

func (m *mockDB) GetAll() []*models.Post {
	var posts []*models.Post
	for _, post := range m.container {
		posts = append(posts, post)
	}
	return posts
}

type TestCase struct {
	posts []*models.Post
	want  []*models.Post
}

func TestPostService(t *testing.T) {
	db := newMockDB()
	service := NewPostServiceImpl(db)

	tests := []TestCase{
		{posts: nil, want: nil},
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
		got := service.GetAll()
		sort.Slice(got, func(i, j int) bool {
			return got[i].Header < got[j].Header
		})
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetAll() = %v, want %v", got, test.want)
		}
		db.ClearAll()
	}
}
