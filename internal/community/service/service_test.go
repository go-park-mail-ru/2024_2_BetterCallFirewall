package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	repo *MockRepo
}

func getService(ctrl *gomock.Controller) (*Service, *mocks) {
	m := &mocks{
		repo: NewMockRepo(ctrl),
	}

	return NewCommunityService(m.repo), m
}

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	res, _ := getService(ctrl)
	assert.NotNil(t, res)
}

var errMock = errors.New("mock error")

func TestGet(t *testing.T) {
	tests := []TableTest[[]*models.CommunityCard, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) ([]*models.CommunityCard, error) {
				return implementation.Get(ctx, input)
			},
			ExpectedResult: func() ([]*models.CommunityCard, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().GetBatch(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) ([]*models.CommunityCard, error) {
				return implementation.Get(ctx, input)
			},
			ExpectedResult: func() ([]*models.CommunityCard, error) {
				return nil, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().GetBatch(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) ([]*models.CommunityCard, error) {
				return implementation.Get(ctx, input)
			},
			ExpectedResult: func() ([]*models.CommunityCard, error) {
				return []*models.CommunityCard{{ID: 1}}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().GetBatch(gomock.Any(), gomock.Any()).Return(
					[]*models.CommunityCard{{ID: 1}},
					nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

func TestGetOne(t *testing.T) {
	tests := []TableTest[*models.Community, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) (*models.Community, error) {
				return implementation.GetOne(ctx, input)
			},
			ExpectedResult: func() (*models.Community, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) (*models.Community, error) {
				return implementation.GetOne(ctx, input)
			},
			ExpectedResult: func() (*models.Community, error) {
				return nil, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) (*models.Community, error) {
				return implementation.GetOne(ctx, input)
			},
			ExpectedResult: func() (*models.Community, error) {
				return &models.Community{ID: 1}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(
					&models.Community{ID: 1},
					nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type InputCreate struct {
	community *models.Community
	authorID  uint32
}

func TestCreate(t *testing.T) {
	tests := []TableTest[struct{}, InputCreate]{
		{
			name: "1",
			SetupInput: func() (*InputCreate, error) {
				input := InputCreate{authorID: 1}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputCreate) (struct{}, error) {
				err := implementation.Create(ctx, input.community, input.authorID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(input InputCreate, m *mocks) {
				m.repo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*InputCreate, error) {
				input := InputCreate{authorID: 1, community: &models.Community{About: "my community"}}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputCreate) (struct{}, error) {
				err := implementation.Create(ctx, input.community, input.authorID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputCreate, m *mocks) {
				m.repo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type InputUpdate struct {
	community *models.Community
	ID        uint32
}

func TestUpdate(t *testing.T) {
	tests := []TableTest[struct{}, InputUpdate]{
		{
			name: "1",
			SetupInput: func() (*InputUpdate, error) {
				input := InputUpdate{ID: 1, community: &models.Community{About: "my community"}}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputUpdate) (struct{}, error) {
				err := implementation.Update(ctx, input.ID, input.community)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(input InputUpdate, m *mocks) {
				m.repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*InputUpdate, error) {
				input := InputUpdate{ID: 1, community: &models.Community{About: "my community"}}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputUpdate) (struct{}, error) {
				err := implementation.Update(ctx, input.ID, input.community)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputUpdate, m *mocks) {
				m.repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []TableTest[struct{}, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) (struct{}, error) {
				err := implementation.Delete(ctx, input)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				input := uint32(1)
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input uint32) (struct{}, error) {
				err := implementation.Delete(ctx, input)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input uint32, m *mocks) {
				m.repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type InputCheckAccess struct {
	userID      uint32
	communityID uint32
}

func TestCheckAccess(t *testing.T) {
	tests := []TableTest[bool, InputCheckAccess]{
		{
			name: "1",
			SetupInput: func() (*InputCheckAccess, error) {
				input := InputCheckAccess{userID: 0, communityID: 0}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputCheckAccess) (bool, error) {
				res := implementation.CheckAccess(ctx, input.userID, input.communityID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return false, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputCheckAccess, m *mocks) {
				m.repo.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
			},
		},
		{
			name: "2",
			SetupInput: func() (*InputCheckAccess, error) {
				input := InputCheckAccess{userID: 1, communityID: 10}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputCheckAccess) (bool, error) {
				res := implementation.CheckAccess(ctx, input.userID, input.communityID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return true, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputCheckAccess, m *mocks) {
				m.repo.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type InputSearch struct {
	query  string
	lastID uint32
}

func TestSearch(t *testing.T) {
	tests := []TableTest[[]*models.CommunityCard, InputSearch]{
		{
			name: "1",
			SetupInput: func() (*InputSearch, error) {
				input := InputSearch{query: "", lastID: 0}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputSearch) ([]*models.CommunityCard, error) {
				res, err := implementation.Search(ctx, input.query, input.lastID)
				return res, err
			},
			ExpectedResult: func() ([]*models.CommunityCard, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(input InputSearch, m *mocks) {
				m.repo.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*InputSearch, error) {
				input := InputSearch{query: "community", lastID: 0}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *Service, input InputSearch) ([]*models.CommunityCard, error) {
				res, err := implementation.Search(ctx, input.query, input.lastID)
				return res, err
			},
			ExpectedResult: func() ([]*models.CommunityCard, error) {
				return []*models.CommunityCard{{ID: 1, Name: "community"}}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputSearch, m *mocks) {
				m.repo.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*models.CommunityCard{{ID: 1, Name: "community"}}, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(*input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, serv, *input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *Service, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
