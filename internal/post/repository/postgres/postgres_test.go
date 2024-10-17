package postgres

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var ID uint32 = 1
	rows := sqlmock.NewRows([]string{"id", "author_id", "content", "created_at", "updated_at"})
	expect := []*models.Post{
		{ID: ID, Header: models.Header{AuthorID: 1}, PostContent: models.Content{Text: "content from user 1", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}
	for _, post := range expect {
		rows = rows.AddRow(post.ID, post.Header.AuthorID, post.PostContent.Text, post.PostContent.CreatedAt, post.PostContent.UpdatedAt)
	}

	mock.ExpectQuery(getPost).WithArgs(1).WillReturnRows(rows)

	repo := NewAdapter(db)

	post, err := repo.Get(context.Background(), ID)
	assert.NoError(t, err, "unexpected err: %s", err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equalf(t, post, expect[0], "results not match, want %v, have %v", expect[0], post)
}
