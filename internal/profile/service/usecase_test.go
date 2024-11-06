package service

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type MockProfileDB struct {
	Storage struct{}
}

type MockPostDB struct {
	Storage struct{}
}

type Test struct {
	ctx              context.Context
	userID           uint32
	profile          *models.FullProfile
	friendID         uint32
	inputProfile     *models.FullProfile
	resProfile       *models.FullProfile
	resShortProfiles []*models.ShortProfile
	resID            []uint32
	resHeader        models.Header

	err error
}

var (
	profileDB = MockProfileDB{}
	postDB    = MockPostDB{}
	pu        = NewProfileUsecase(profileDB, postDB)

	examplePost = &models.Post{
		ID:          1,
		Header:      models.Header{AuthorID: 1, Author: "Andrew"},
		PostContent: models.Content{Text: "Hello, World!"},
	}

	exampleProfileWithPost = &models.FullProfile{
		ID:        1,
		FirstName: "Andrew",
		LastName:  "Savvateev",
		Bio:       "Hello, viewers!",
		Avatar:    "/default",
		Pics:      nil,
		Posts:     []*models.Post{examplePost},
	}
	shortExample1 = &models.ShortProfile{
		ID:        1,
		FirstName: "Andrew",
		LastName:  "Savvateev",
		Avatar:    "",
	}

	exampleProfileWithoutPost = &models.FullProfile{
		ID:        2,
		FirstName: "Alex",
		LastName:  "Zem",
		Bio:       "Hello, viewers!",
		Avatar:    "",
		Pics:      nil,
		Posts:     nil,
	}
	shortExample2 = &models.ShortProfile{
		ID:        2,
		FirstName: "Alex",
		LastName:  "Zem",
		Avatar:    "",
	}
)

var ErrExec = errors.New("execution error")

func (m MockProfileDB) GetProfileById(ctx context.Context, u uint32) (*models.FullProfile, error) {
	if u == 1 {
		return exampleProfileWithPost, nil
	}
	if u == 2 {
		return exampleProfileWithoutPost, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) GetAll(ctx context.Context, self uint32, lastId uint32) ([]*models.ShortProfile, error) {
	if self == 3 {
		return []*models.ShortProfile{shortExample1, shortExample2}, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) UpdateProfile(ctx context.Context, profile *models.FullProfile) error {
	return nil
}

func (m MockProfileDB) DeleteProfile(u uint32) error {
	return nil
}

func (m MockProfileDB) AddFriendsReq(reciever uint32, sender uint32) error {
	if reciever == 10 || sender == 10 {
		return ErrExec
	}
	return nil
}

func (m MockProfileDB) AcceptFriendsReq(who uint32, whose uint32) error {
	if who == 10 || whose == 10 {
		return ErrExec
	}
	return nil
}

func (m MockProfileDB) MoveToSubs(who uint32, whom uint32) error {
	if who == 10 || whom == 10 {
		return ErrExec
	}
	return nil
}

func (m MockProfileDB) RemoveSub(who uint32, whom uint32) error {
	if who == 10 || whom == 10 {
		return ErrExec
	}
	return nil
}

func (m MockProfileDB) GetAllFriends(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error) {
	if u == 3 {
		return []*models.ShortProfile{shortExample1, shortExample2}, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) GetAllSubs(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error) {
	if u == 3 {
		return []*models.ShortProfile{shortExample1, shortExample2}, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) GetAllSubscriptions(ctx context.Context, u uint32, u2 uint32) ([]*models.ShortProfile, error) {
	if u == 3 {
		return []*models.ShortProfile{shortExample1, shortExample2}, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) GetFriendsID(ctx context.Context, u uint32) ([]uint32, error) {
	if u == 3 {
		return []uint32{1, 2}, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) GetHeader(ctx context.Context, u uint32) (*models.Header, error) {
	if u == 1 {
		return &models.Header{AuthorID: u, Author: "Andrew Savvateev"}, nil
	}
	return nil, sql.ErrNoRows
}

func (m MockProfileDB) GetStatus(context.Context, uint32, uint32) (int, error) { return 0, nil }

func (m MockProfileDB) UpdateWithAvatar(context.Context, *models.FullProfile) error { return nil }

func (m MockProfileDB) GetSubscriptionsID(context.Context, uint32) ([]uint32, error) { return nil, nil }

func (m MockProfileDB) GetSubscribersID(context.Context, uint32) ([]uint32, error) { return nil, nil }
func (m MockProfileDB) GetStatuses(context.Context, uint32) ([]uint32, []uint32, []uint32, error) {
	return nil, nil, nil, nil
}

func (m MockPostDB) GetAuthorsPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	if header.AuthorID == 1 {
		return []*models.Post{examplePost}, nil
	}
	/*if header.AuthorID == 2 {
		return nil, nil
	}*/
	return nil, nil
}

func TestGetProfileByID(t *testing.T) {
	sessId1, err := models.NewSession(1)
	if err != nil {
		t.Fatal(err)
	}

	sessId10, err := models.NewSession(10)
	if err != nil {
		t.Fatal(err)
	}

	sessId2, err := models.NewSession(2)
	if err != nil {
		t.Fatal(err)
	}
	tests := []Test{
		{
			ctx:        models.ContextWithSession(context.Background(), sessId1),
			userID:     1,
			resProfile: exampleProfileWithPost,
			err:        nil,
		},
		{
			ctx:    models.ContextWithSession(context.Background(), sessId10),
			userID: 10,
			err:    sql.ErrNoRows,
		},
		{
			ctx:        models.ContextWithSession(context.Background(), sessId2),
			userID:     2,
			resProfile: exampleProfileWithoutPost,
			err:        nil,
		},
	}

	for caseNum, test := range tests {
		res, err := pu.GetProfileById(test.ctx, test.userID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.resProfile) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.resProfile, res)
		}
	}
}

func TestGetAll(t *testing.T) {
	tests := []Test{
		{
			ctx:              context.Background(),
			userID:           3,
			resShortProfiles: []*models.ShortProfile{shortExample1, shortExample2},
			err:              nil,
		},
		{
			ctx:    context.Background(),
			userID: 10,
			err:    sql.ErrNoRows,
		},
	}

	for caseNum, test := range tests {
		res, err := pu.GetAll(test.ctx, test.userID, 0)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.resShortProfiles) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.resShortProfiles, res)
		}
	}
}

func TestUpdateProfile(t *testing.T) {
	tests := []Test{
		{
			ctx:          context.Background(),
			inputProfile: exampleProfileWithPost,
			err:          nil,
		},
	}

	for caseNum, test := range tests {
		err := pu.UpdateProfile(nil, test.inputProfile)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}

func TestDeleteProfile(t *testing.T) {
	tests := []Test{
		{
			ctx:          context.Background(),
			userID:       1,
			inputProfile: exampleProfileWithPost,
			err:          nil,
		},
	}

	for caseNum, test := range tests {
		err := pu.DeleteProfile(test.userID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}

func TestSendFriendReq(t *testing.T) {
	tests := []Test{
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 2,
			err:      nil,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 1,
			err:      myErr.ErrSameUser,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 10,
			err:      ErrExec,
		},
	}

	for caseNum, test := range tests {
		err := pu.SendFriendReq(test.userID, test.friendID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}

func TestAcceptFriendReq(t *testing.T) {
	tests := []Test{
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 2,
			err:      nil,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 1,
			err:      myErr.ErrSameUser,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 10,
			err:      ErrExec,
		},
	}

	for caseNum, test := range tests {
		err := pu.AcceptFriendReq(test.userID, test.friendID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}

func TestRemoveFromFriends(t *testing.T) {
	tests := []Test{
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 2,
			err:      nil,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 1,
			err:      myErr.ErrSameUser,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 10,
			err:      ErrExec,
		},
	}

	for caseNum, test := range tests {
		err := pu.RemoveFromFriends(test.userID, test.friendID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	tests := []Test{
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 2,
			err:      nil,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 1,
			err:      myErr.ErrSameUser,
		},
		{
			ctx:      context.Background(),
			userID:   1,
			friendID: 10,
			err:      ErrExec,
		},
	}

	for caseNum, test := range tests {
		err := pu.Unsubscribe(test.userID, test.friendID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}

func TestGetAllFriends(t *testing.T) {
	sessId3, err := models.NewSession(3)
	if err != nil {
		t.Fatal(err)
	}

	sessId10, err := models.NewSession(10)
	if err != nil {
		t.Fatal(err)
	}
	tests := []Test{
		{
			ctx:              models.ContextWithSession(context.Background(), sessId3),
			userID:           3,
			resShortProfiles: []*models.ShortProfile{shortExample1, shortExample2},
			err:              nil,
		},
		{
			ctx:    models.ContextWithSession(context.Background(), sessId10),
			userID: 10,
			err:    sql.ErrNoRows,
		},
	}

	for caseNum, test := range tests {
		res, err := pu.GetAllFriends(test.ctx, test.userID, 0)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.resShortProfiles) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.resShortProfiles, res)
		}
	}
}

func TestGetAllSubs(t *testing.T) {
	sessId3, err := models.NewSession(3)
	if err != nil {
		t.Fatal(err)
	}

	sessId10, err := models.NewSession(10)
	if err != nil {
		t.Fatal(err)
	}
	tests := []Test{
		{
			ctx:              models.ContextWithSession(context.Background(), sessId3),
			userID:           3,
			resShortProfiles: []*models.ShortProfile{shortExample1, shortExample2},
			err:              nil,
		},
		{
			ctx:    models.ContextWithSession(context.Background(), sessId10),
			userID: 10,
			err:    sql.ErrNoRows,
		},
	}

	for caseNum, test := range tests {
		res, err := pu.GetAllSubs(test.ctx, test.userID, 0)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.resShortProfiles) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.resShortProfiles, res)
		}
	}
}

func TestGetAllSubscriptions(t *testing.T) {
	sessId3, err := models.NewSession(3)
	if err != nil {
		t.Fatal(err)
	}

	sessId10, err := models.NewSession(10)
	if err != nil {
		t.Fatal(err)
	}
	tests := []Test{
		{
			ctx:              models.ContextWithSession(context.Background(), sessId3),
			userID:           3,
			resShortProfiles: []*models.ShortProfile{shortExample1, shortExample2},
			err:              nil,
		},
		{
			ctx:    models.ContextWithSession(context.Background(), sessId10),
			userID: 10,
			err:    sql.ErrNoRows,
		},
	}

	for caseNum, test := range tests {
		res, err := pu.GetAllSubscriptions(test.ctx, test.userID, 0)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.resShortProfiles) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.resShortProfiles, res)
		}
	}
}

func TestGetHeader(t *testing.T) {
	tests := []Test{
		{
			ctx:       context.Background(),
			userID:    1,
			resHeader: models.Header{AuthorID: 1, Author: "Andrew Savvateev"},
			err:       nil,
		},
		{
			ctx:    context.Background(),
			userID: 10,
			err:    sql.ErrNoRows,
		},
	}

	for caseNum, test := range tests {
		res, err := pu.GetHeader(test.ctx, test.userID)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.resHeader) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.resShortProfiles, res)
		}
	}
}
