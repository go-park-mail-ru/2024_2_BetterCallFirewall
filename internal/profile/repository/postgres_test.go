package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type Test struct {
	inputID      uint32
	friendID     uint32
	inputProfile *models.FullProfile
	execResult   driver.Result
	resProfile   *models.FullProfile
	resProfiles  []*models.ShortProfile
	resHeader    *models.Header
	resIDs       []uint32
	expectedErr  error
	dbError      error
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
			inputID:     0,
			resProfiles: nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID:     3,
			resProfiles: nil,
			expectedErr: errors.New("sql: expected 2 destination arguments in Scan, not 4"),
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 2 {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllProfilesBatch)).
				WithArgs(test.inputID, uint32(0), 20).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllProfilesBatch)).
				WithArgs(test.inputID, 0, 20).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}
		profiles, err := ProfileManager.GetAll(context.Background(), test.inputID, 0)

		if casenum == 2 {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
			}
			assert.Equalf(t, test.resProfiles, profiles, "case [%d]: results must match, want %v, have %v", casenum, test.resProfile, profiles)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
			}
		}
	}
}

func TestUpdateProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	updatedProfile := &models.FullProfile{
		ID:        1,
		FirstName: "Andrew",
		LastName:  "Savvateev",
		Bio:       "Hello, world",
		Avatar:    "/default",
		Pics:      nil,
		Posts:     nil,
	}

	tests := []Test{
		{
			inputProfile: updatedProfile,
			execResult:   sqlmock.NewResult(1, 3),
			expectedErr:  nil,
			dbError:      nil,
		},
		{
			inputProfile: &models.FullProfile{ID: 10},
			expectedErr:  errMockDb,
			dbError:      errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(UpdateProfile)).
			WithArgs(test.inputProfile.FirstName, test.inputProfile.LastName, test.inputProfile.Bio, test.inputProfile.ID).
			WillReturnResult(test.execResult).
			WillReturnError(test.dbError)
		err := ProfileManager.UpdateProfile(context.Background(), test.inputProfile)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestDeleteProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profileId uint32 = 1

	tests := []Test{
		{
			inputID:     profileId,
			execResult:  sqlmock.NewResult(1, 7),
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     10,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(DeleteProfile)).
			WithArgs(test.inputID).
			WillReturnResult(test.execResult).
			WillReturnError(test.dbError)
		err := ProfileManager.DeleteProfile(test.inputID)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestAddFriendsReq(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profileId uint32 = 1
	var friend uint32 = 2

	tests := []Test{
		{
			inputID:     profileId,
			friendID:    friend,
			execResult:  sqlmock.NewResult(1, 3),
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     profileId,
			friendID:    27,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(AddFriends)).
			WithArgs(test.inputID, test.friendID).
			WillReturnResult(test.execResult).
			WillReturnError(test.dbError)
		err := ProfileManager.AddFriendsReq(test.friendID, test.inputID)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestAcceptFriendsReq(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profileId uint32 = 1
	var friend uint32 = 2

	tests := []Test{
		{
			inputID:     profileId,
			friendID:    friend,
			execResult:  sqlmock.NewResult(1, 3),
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     profileId,
			friendID:    27,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(AcceptFriendReq)).
			WithArgs(test.friendID, test.inputID).
			WillReturnResult(test.execResult).
			WillReturnError(test.dbError)
		err := ProfileManager.AcceptFriendsReq(test.inputID, test.friendID)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestMoveToSubs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profileId uint32 = 1
	var friend uint32 = 2

	tests := []Test{
		{
			inputID:     profileId,
			friendID:    friend,
			execResult:  sqlmock.NewResult(1, 3),
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     profileId,
			friendID:    27,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(RemoveFriendsReq)).
			WithArgs(test.inputID, test.friendID).
			WillReturnResult(test.execResult).
			WillReturnError(test.dbError)
		err := ProfileManager.MoveToSubs(test.inputID, test.friendID)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestRemoveSub(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profileId uint32 = 1
	var friend uint32 = 2

	tests := []Test{
		{
			inputID:     profileId,
			friendID:    friend,
			execResult:  sqlmock.NewResult(1, 3),
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     profileId,
			friendID:    27,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		mock.ExpectExec(regexp.QuoteMeta(DeleteFriendship)).
			WithArgs(test.inputID, test.friendID).
			WillReturnResult(test.execResult).
			WillReturnError(test.dbError)
		err := ProfileManager.RemoveSub(test.inputID, test.friendID)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
		}
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
		}
	}
}

func TestGetAllFriends(t *testing.T) {
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
			inputID:     0,
			resProfiles: nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID:     3,
			resProfiles: nil,
			expectedErr: errors.New("sql: expected 2 destination arguments in Scan, not 4"),
			dbError:     errMockDb,
		},
	}

	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 2 {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllFriends)).
				WithArgs(test.inputID, 0, LIMIT).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllFriends)).
				WithArgs(test.inputID, 0, LIMIT).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}

		profiles, err := ProfileManager.GetAllFriends(context.Background(), test.inputID, 0)
		if casenum == 2 {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
			}
			assert.Equalf(t, test.resProfiles, profiles, "case [%d]: results must match, want %v, have %v", casenum, test.resProfile, profiles)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
			}
		}
	}
}

func TestGetAllSubs(t *testing.T) {
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
			inputID:     0,
			resProfiles: nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID:     3,
			resProfiles: nil,
			expectedErr: errors.New("sql: expected 2 destination arguments in Scan, not 4"),
			dbError:     errMockDb,
		},
	}
	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 2 {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllSubs)).
				WithArgs(test.inputID, 0, LIMIT).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllSubs)).
				WithArgs(test.inputID, 0, LIMIT).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}

		profiles, err := ProfileManager.GetAllSubs(context.Background(), test.inputID, 0)
		if casenum == 2 {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
			}
			assert.Equalf(t, test.resProfiles, profiles, "case [%d]: results must match, want %v, have %v", casenum, test.resProfile, profiles)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
			}
		}
	}
}

func TestGetAllSubscriptions(t *testing.T) {
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
			inputID:     0,
			resProfiles: nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID:     3,
			resProfiles: nil,
			expectedErr: errors.New("sql: expected 2 destination arguments in Scan, not 4"),
			dbError:     errMockDb,
		},
	}
	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 2 {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllSubscriptions)).
				WithArgs(test.inputID, 0, LIMIT).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetAllSubscriptions)).
				WithArgs(test.inputID, 0, LIMIT).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}

		profiles, err := ProfileManager.GetAllSubscriptions(context.Background(), test.inputID, 0)
		if casenum == 2 {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
			}
			assert.Equalf(t, test.resProfiles, profiles, "case [%d]: results must match, want %v, have %v", casenum, test.resProfile, profiles)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
			}
		}
	}
}

func TestGetFriendsID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profID uint32 = 3

	rows := sqlmock.NewRows([]string{"id"})
	expect := []uint32{1, 2, 3}
	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	errRows := sqlmock.NewRows([]string{"id", "first_name"}).AddRow(3, "Andrew")
	tests := []Test{
		{
			inputID:     profID,
			resIDs:      expect,
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     0,
			resIDs:      nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID: 3,
			resIDs:  nil,
		},
	}
	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 2 {
			mock.ExpectQuery(regexp.QuoteMeta(GetFriendsID)).
				WithArgs(test.inputID).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetFriendsID)).
				WithArgs(test.inputID).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}

		profiles, err := ProfileManager.GetFriendsID(context.Background(), test.inputID)
		if casenum == 2 {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
			}
			assert.Equalf(t, test.resIDs, profiles, "case [%d]: results must match, want %v, have %v", casenum, test.resIDs, profiles)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
			}
		}
	}
}

func TestGetHeader(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var profID uint32 = 3

	rows := sqlmock.NewRows([]string{"name", "avatar"})
	expect := []*models.Header{
		{
			Author: "Andrew Savvateev",
			Avatar: "/default",
		},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.Author, item.Avatar)
	}

	errRows := sqlmock.NewRows([]string{"id", "first_name", "last_name"}).AddRow(3, "Andrew", "Savvateev")
	tests := []Test{
		{
			inputID: profID,
			resHeader: &models.Header{
				AuthorID: profID,
				Author:   "Andrew Savvateev",
				Avatar:   "/default",
			},
			expectedErr: nil,
			dbError:     nil,
		},
		{
			inputID:     0,
			resHeader:   nil,
			expectedErr: errMockDb,
			dbError:     errMockDb,
		},
		{
			inputID:   3,
			resHeader: nil,
		},
	}
	ProfileManager := NewProfileRepo(db)
	for casenum, test := range tests {
		if casenum == 2 {
			mock.ExpectQuery(regexp.QuoteMeta(GetShortProfile)).
				WithArgs(test.inputID).WillReturnRows(errRows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(GetShortProfile)).
				WithArgs(test.inputID).WillReturnRows(rows).
				WillReturnError(test.dbError)
		}

		profile, err := ProfileManager.GetHeader(context.Background(), test.inputID)
		if casenum == 2 {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		} else {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("case [%d]: there were unfulfilled expectations: %v", casenum, err)
			}
			assert.Equalf(t, test.resHeader, profile, "case [%d]: results must match, want %v, have %v", casenum, test.resIDs, profile)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("case [%d]: errors must match, have %v, want %v", casenum, err, test.expectedErr)
			}
		}
	}
}
