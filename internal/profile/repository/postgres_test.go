package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type Test struct {
	inputID     uint32
	resProfile  *models.FullProfile
	resProfiles []*models.ShortProfile
	expectedErr error
	dbError     error
}

var errMockDb = errors.New("mock db error")

func TestGetProfileById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profID uint32 = 1

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "bio", "avatar"})
	expect := []*models.FullProfile{
		{
			ID:        profID,
			FirstName: "Andrew",
			LastName:  "Savvateev",
			Bio:       "Hello, world",
			Avatar:    "/default",
			Pics:      nil,
			Posts:     nil,
		},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FirstName, item.LastName, item.Bio, item.Avatar)
	}

	tests := []Test{
		{
			inputID:     profID,
			resProfile:  expect[0],
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     100,
			resProfile:  nil,
			expectedErr: myErr.ErrProfileNotFound,
			dbError:     sql.ErrNoRows,
		},
		{
			inputID:     0,
			resProfile:  nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectQuery(regexp.QuoteMeta(GetProfileByID)).
			WithArgs(test.inputID).WillReturnRows(rows).
			WillReturnError(test.dbError)
		profile, err := ProfileManager.GetProfileById(context.Background(), test.inputID)
		assert.Equalf(t, test.resProfile, profile, "case [%d]: results must match, want %v, have %v", casenum, test.resProfile, profile)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profID uint32 = 3

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "avatar"})
	expect := []*models.ShortProfile{
		{
			ID:        1,
			FirstName: "Andrew",
			LastName:  "Savvateev",
			Avatar:    "/default",
		},
		{
			ID:        2,
			FirstName: "Alex",
			LastName:  "Zem",
			Avatar:    "/default",
		},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.FirstName, item.LastName, item.Avatar)
	}

	errRows := sqlmock.NewRows([]string{"id", "first_name"}).AddRow(3, "Andrew")
	tests := []Test{
		{
			inputID:     profID,
			resProfiles: []*models.ShortProfile{expect[0], expect[1]},
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     1,
			resProfiles: []*models.ShortProfile{expect[1]},
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     0,
			resProfiles: nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID:     3,
			resProfiles: nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 3 {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllProfiles)).
				WithArgs(test.inputID).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllProfiles)).
				WithArgs(test.inputID).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}
		profiles, err := ProfileManager.GetProfileById(context.Background(), test.inputID)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
		assert.Equalf(t, test.resProfile, profiles, "case [%d]: results must match, want %v, have %v", casenum, test.resProfile, profiles)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
	}
}

func TestUpdateProfile(t *testing.T) {

}

func TestDeleteProfile(t *testing.T) {

}

func TestAddFriendsReq(t *testing.T) {

}

func TestAcceptFriendsReq(t *testing.T) {

}

func TestMoveToSubs(t *testing.T) {

}

func TestRemoveSub(t *testing.T) {

}

func TestGetAllFriends(t *testing.T) {

}

func TestGetAllSubs(t *testing.T) {

}

func TestGetAllSubscriptions(t *testing.T) {

}

func TestGetFriendsID(t *testing.T) {

}

func TestGetHeader(t *testing.T) {

}
