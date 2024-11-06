package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

var errMockDB = errors.New("mock db error")

type TestCaseGet struct {
	ID       uint32
	wantErr  error
	dbErr    error
	wantPost *models.Post
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var ID uint32 = 1
	rows := sqlmock.NewRows([]string{"id", "author_id", "content", "created_at"})
	expect := []*models.Post{
		{ID: ID, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: time.Now()}},
	}
	for _, post := range expect {
		rows = rows.AddRow(ID, post.Header.AuthorID, post.PostContent.Text, post.PostContent.CreatedAt)
	}

	repo := NewAdapter(db)

	tests := []TestCaseGet{
		{ID: ID, wantPost: expect[0], wantErr: nil, dbErr: nil},
		{ID: 100, wantPost: nil, wantErr: myErr.ErrPostNotFound, dbErr: sql.ErrNoRows},
		{ID: 10, wantPost: nil, wantErr: errMockDB, dbErr: errMockDB},
	}

	for _, test := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(getPost)).
			WithArgs(test.ID).
			WillReturnRows(rows).
			WillReturnError(test.dbErr)

		post, err := repo.Get(context.Background(), test.ID)
		assert.Equalf(t, test.wantPost, post, "results not match,\n want %v\n have %v", test.wantPost, post)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("errors not match,\n want %v\n have %v", test.wantErr, err)
		}
	}
}

type TestCaseCreate struct {
	post    *models.Post
	wantID  uint32
	wantErr error
	dbErr   error
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAdapter(db)

	tests := []TestCaseCreate{
		{post: &models.Post{Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1"}}, wantID: 1, wantErr: nil, dbErr: nil},
		{post: &models.Post{Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2"}}, wantID: 2, wantErr: nil, dbErr: nil},
		{post: &models.Post{Header: models.Header{AuthorID: 10}, PostContent: models.Content{Text: "wrong query"}}, wantID: 0, wantErr: errMockDB, dbErr: errMockDB},
	}

	for _, test := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(createPost)).
			WithArgs(test.post.Header.AuthorID, test.post.PostContent.Text).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.wantID)).
			WillReturnError(test.dbErr)

		id, err := repo.Create(context.Background(), test.post)
		if id != test.wantID {
			t.Errorf("results not match,\n want %v\n have %v", test.wantID, id)
		}
		if !errors.Is(err, test.wantErr) {
			t.Errorf("unexpected err:\n want:%v\n got:%v", test.wantErr, err)
		}
	}
}

type TestCaseDelete struct {
	ID           uint32
	wantErr      error
	dbErr        error
	rowsAffected int64
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAdapter(db)

	tests := []TestCaseDelete{
		{ID: 1, wantErr: nil, rowsAffected: 1, dbErr: nil},
		{ID: 1, wantErr: myErr.ErrPostNotFound, rowsAffected: 0, dbErr: nil},
		{ID: 100, wantErr: myErr.ErrPostNotFound, rowsAffected: 0, dbErr: nil},
		{ID: 10, wantErr: errMockDB, rowsAffected: 0, dbErr: errMockDB},
	}

	for _, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(deletePost)).
			WithArgs(test.ID).
			WillReturnResult(sqlmock.NewResult(0, test.rowsAffected)).
			WillReturnError(test.dbErr)

		err := repo.Delete(context.Background(), test.ID)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("errors not match,\n want %v\n have %v", test.wantErr, err)
		}
	}
}

type TestCaseUpdate struct {
	post         *models.Post
	rowsAffected int64
	wantErr      error
	dbErr        error
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewAdapter(db)

	tests := []TestCaseUpdate{
		{post: &models.Post{ID: 1, PostContent: models.Content{Text: "update post", UpdatedAt: time.Now()}}, wantErr: nil, dbErr: nil, rowsAffected: 1},
		{post: &models.Post{ID: 2, PostContent: models.Content{Text: "wrong ID", UpdatedAt: time.Now()}}, wantErr: myErr.ErrPostNotFound, dbErr: nil, rowsAffected: 0},
		{post: &models.Post{ID: 1, PostContent: models.Content{Text: "update post who was update early", UpdatedAt: time.Now()}}, wantErr: nil, dbErr: nil, rowsAffected: 1},
		{post: &models.Post{ID: 5, PostContent: models.Content{Text: "wrong query", UpdatedAt: time.Now()}}, wantErr: errMockDB, dbErr: errMockDB, rowsAffected: 0},
	}

	for _, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(updatePost)).
			WithArgs(test.post.PostContent.Text, test.post.PostContent.UpdatedAt, test.post.ID).
			WillReturnResult(sqlmock.NewResult(0, test.rowsAffected)).
			WillReturnError(test.dbErr)

		err := repo.Update(context.Background(), test.post)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("errors not match,\n want %v\n have %v", test.wantErr, err)
		}
	}
}

type TestCaseGetPostAuthor struct {
	postID       uint32
	wantAuthorID uint32
	wantErr      error
	dbErr        error
}

func TestGetPostAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var ID uint32 = 1

	repo := NewAdapter(db)

	tests := []TestCaseGetPostAuthor{
		{postID: ID, wantAuthorID: 1, wantErr: nil},
		{postID: 100, wantAuthorID: 0, wantErr: myErr.ErrPostNotFound, dbErr: sql.ErrNoRows},
		{postID: 100, wantAuthorID: 0, wantErr: errMockDB, dbErr: errMockDB},
	}

	for _, test := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(getPostAuthor)).
			WithArgs(test.postID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.wantAuthorID)).
			WillReturnError(test.dbErr)

		id, err := repo.GetPostAuthor(context.Background(), test.postID)
		assert.Equalf(t, test.wantAuthorID, id, "results not match,\n want %v\n have %v", test.wantAuthorID, id)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("errors not match,\n want %v\n have %v", test.wantErr, err)
		}
	}
}

type TestCaseGetAuthorPosts struct {
	Author    *models.Header
	wantPosts []*models.Post
	dbErr     error
	wantErr   error
}

func TestGetAuthorsPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	createTime := time.Now()

	repo := NewAdapter(db)
	tests := []TestCaseGetAuthorPosts{
		{
			Author: &models.Header{AuthorID: 1},
			wantPosts: []*models.Post{
				{ID: 1, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
				{ID: 2, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "another content from user 1", CreatedAt: createTime}},
			},
			wantErr: nil,
			dbErr:   nil,
		},
		{
			Author: &models.Header{AuthorID: 2},
			wantPosts: []*models.Post{
				{ID: 3, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
			},
			wantErr: nil,
			dbErr:   nil,
		},
		{
			Author:    &models.Header{AuthorID: 3},
			wantPosts: nil,
			wantErr:   myErr.ErrNoMoreContent,
			dbErr:     sql.ErrNoRows,
		},
		{
			Author:    &models.Header{AuthorID: 4},
			wantPosts: nil,
			wantErr:   errMockDB,
			dbErr:     errMockDB,
		},
	}

	for _, test := range tests {
		rows := sqlmock.NewRows([]string{"id", "content", "created_at"})
		for _, post := range test.wantPosts {
			rows = rows.AddRow(post.ID, post.PostContent.Text, post.PostContent.CreatedAt)
		}
		mock.ExpectQuery(regexp.QuoteMeta(getProfilePosts)).
			WithArgs(test.Author.AuthorID).
			WillReturnRows(rows).
			WillReturnError(test.dbErr)

		posts, err := repo.GetAuthorPosts(context.Background(), test.Author)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("unexpected error: got:%v\nwant:%v\n", err, test.wantErr)
		}
		assert.Equalf(t, posts, test.wantPosts, "result dont match\nwant: %v\ngot:%v", test.wantPosts, posts)
	}
}

type TestCaseGetPosts struct {
	lastID   uint32
	wantPost []*models.Post
	dbErr    error
	wantErr  error
}

func TestGetPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	createTime := time.Now()
	expect := []*models.Post{
		{ID: 1, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
		{ID: 2, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
		{ID: 3, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
		{ID: 4, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
		{ID: 5, Header: models.Header{AuthorID: 3}, PostContent: models.Content{Text: "content from user 3", CreatedAt: createTime}},
		{ID: 6, Header: models.Header{AuthorID: 3}, PostContent: models.Content{Text: "content from user 3", CreatedAt: createTime}},
		{ID: 7, Header: models.Header{AuthorID: 6}, PostContent: models.Content{Text: "content from user 6", CreatedAt: createTime}},
		{ID: 8, Header: models.Header{AuthorID: 4}, PostContent: models.Content{Text: "content from user 4", CreatedAt: createTime}},
		{ID: 9, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
		{ID: 10, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
		{ID: 11, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
	}

	repo := NewAdapter(db)

	tests := []TestCaseGetPosts{
		{lastID: 0, wantPost: nil, wantErr: myErr.ErrNoMoreContent, dbErr: sql.ErrNoRows},
		{lastID: 1, wantPost: nil, wantErr: errMockDB, dbErr: errMockDB},
		{
			lastID:   3,
			wantPost: expect[:3],
			wantErr:  nil,
			dbErr:    nil,
		},
		{
			lastID:   11,
			wantPost: expect[1:11],
			wantErr:  nil,
			dbErr:    nil,
		},
	}

	for _, test := range tests {
		rows := sqlmock.NewRows([]string{"id", "author_id", "content", "created_at"})
		for _, post := range test.wantPost {
			rows.AddRow(post.ID, post.Header.AuthorID, post.PostContent.Text, post.PostContent.CreatedAt)
		}
		mock.ExpectQuery(regexp.QuoteMeta(getPostBatch)).
			WithArgs(test.lastID).
			WillReturnRows(rows).
			WillReturnError(test.dbErr)

		posts, err := repo.GetPosts(context.Background(), test.lastID)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("unexpected error: got:%v\nwant:%v\n", err, test.wantErr)
		}
		assert.Equalf(t, posts, test.wantPost, "result dont match\nwant: %v\ngot:%v", test.wantPost, posts)
	}
}

type TestCaseConvertSliceToString struct {
	ids  []uint32
	want string
}

func TestConvertSliceToString(t *testing.T) {
	tests := []TestCaseConvertSliceToString{
		{ids: []uint32{1, 2, 3}, want: "{1, 2, 3}"},
		{ids: nil, want: "{}"},
		{ids: []uint32{2, 3, 5, 20, 30, 50, 60, 9, 8}, want: "{2, 3, 5, 20, 30, 50, 60, 9, 8}"},
	}

	for _, test := range tests {
		res := convertSliceToString(test.ids)
		if res != test.want {
			t.Errorf("unexpected result: got:%v\nwant:%v\n", res, test.want)
		}
	}
}

type GetFriendsPosts struct {
	lastID    uint32
	friendsID []uint32
	wantPost  []*models.Post
	wantErr   error
	dbErr     error
}

func TestGetFriendsPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	createTime := time.Now()

	expect := []*models.Post{
		{ID: 1, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
		{ID: 2, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
		{ID: 3, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
		{ID: 4, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
		{ID: 5, Header: models.Header{AuthorID: 3}, PostContent: models.Content{Text: "content from user 3", CreatedAt: createTime}},
		{ID: 6, Header: models.Header{AuthorID: 3}, PostContent: models.Content{Text: "content from user 3", CreatedAt: createTime}},
		{ID: 7, Header: models.Header{AuthorID: 6}, PostContent: models.Content{Text: "content from user 6", CreatedAt: createTime}},
		{ID: 8, Header: models.Header{AuthorID: 4}, PostContent: models.Content{Text: "content from user 4", CreatedAt: createTime}},
		{ID: 9, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
		{ID: 10, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: createTime}},
		{ID: 11, Header: models.Header{AuthorID: 2}, PostContent: models.Content{Text: "content from user 2", CreatedAt: createTime}},
	}

	repo := NewAdapter(db)

	tests := []GetFriendsPosts{
		{lastID: 0, friendsID: []uint32{}, wantPost: nil, wantErr: myErr.ErrNoMoreContent, dbErr: sql.ErrNoRows},
		{lastID: 1, friendsID: []uint32{}, wantPost: nil, wantErr: errMockDB, dbErr: errMockDB},
		{lastID: 3, friendsID: []uint32{1, 2}, wantPost: expect[:3], wantErr: nil, dbErr: nil},
		{lastID: 11, friendsID: []uint32{3, 6, 4}, wantPost: expect[4:7], wantErr: nil, dbErr: nil},
		{lastID: 11, friendsID: []uint32{1, 2, 5, 3, 6, 4}, wantPost: expect[1:], wantErr: nil, dbErr: nil},
		{lastID: 11, friendsID: []uint32{}, wantPost: nil, wantErr: myErr.ErrNoMoreContent, dbErr: nil},
	}

	for _, test := range tests {
		rows := sqlmock.NewRows([]string{"id", "author_id", "content", "created_at"})
		for _, post := range test.wantPost {
			rows.AddRow(post.ID, post.Header.AuthorID, post.PostContent.Text, post.PostContent.CreatedAt)
		}
		mock.ExpectQuery(regexp.QuoteMeta(getFriendsPost)).
			WithArgs(test.lastID, convertSliceToString(test.friendsID)).
			WillReturnRows(rows).
			WillReturnError(test.dbErr)

		posts, err := repo.GetFriendsPosts(context.Background(), test.friendsID, test.lastID)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("unexpected error: got:%v\nwant:%v\n", err, test.wantErr)
		}
		assert.Equalf(t, posts, test.wantPost, "result dont match\nwant: %v\ngot:%v", test.wantPost, posts)
	}
}
