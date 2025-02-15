package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var (
	errMock = errors.New("mock error")
	pic     = models.Picture("/image/sticker")
)

func getUseCase(ctrl *gomock.Controller) (*StickerUsecaseImplementation, *mocks) {
	m := &mocks{
		repository: NewMockRepository(ctrl),
	}

	return NewStickerUsecase(m.repository), m
}

type input struct {
	userID   uint32
	filepath string
}

func TestAddNewSticker(t *testing.T) {
	tests := []TableTest[struct{}, input]{
		{
			name: "1",
			SetupInput: func() (*input, error) {
				return &input{filepath: "", userID: 0}, nil
			},
			Run: func(
				ctx context.Context, implementation *StickerUsecaseImplementation, request input,
			) (struct{}, error) {
				err := implementation.AddNewSticker(ctx, request.filepath, request.userID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request input, m *mocks) {
				m.repository.EXPECT().AddNewSticker(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*input, error) {
				return &input{filepath: "/image/sticker", userID: 1}, nil
			},
			Run: func(
				ctx context.Context, implementation *StickerUsecaseImplementation, request input,
			) (struct{}, error) {
				err := implementation.AddNewSticker(ctx, request.filepath, request.userID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request input, m *mocks) {
				m.repository.EXPECT().AddNewSticker(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getUseCase(ctrl)
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
			},
		)
	}
}

func TestGetMineSticker(t *testing.T) {
	tests := []TableTest[[]*models.Picture, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				req := uint32(0)
				return &req, nil
			},
			Run: func(
				ctx context.Context, implementation *StickerUsecaseImplementation, request uint32,
			) ([]*models.Picture, error) {
				return implementation.GetMineStickers(ctx, request)
			},
			ExpectedResult: func() ([]*models.Picture, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request uint32, m *mocks) {
				m.repository.EXPECT().GetMineStickers(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				req := uint32(10)
				return &req, nil
			},
			Run: func(
				ctx context.Context, implementation *StickerUsecaseImplementation, request uint32,
			) ([]*models.Picture, error) {
				return implementation.GetMineStickers(ctx, request)
			},
			ExpectedResult: func() ([]*models.Picture, error) {
				return []*models.Picture{&pic}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request uint32, m *mocks) {
				m.repository.EXPECT().GetMineStickers(gomock.Any(), gomock.Any()).Return([]*models.Picture{&pic}, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getUseCase(ctrl)
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
			},
		)
	}
}

func TestGetAllSticker(t *testing.T) {
	tests := []TableTest[[]*models.Picture, struct{}]{
		{
			name: "1",
			SetupInput: func() (*struct{}, error) {
				return &struct{}{}, nil
			},
			Run: func(
				ctx context.Context, implementation *StickerUsecaseImplementation, request struct{},
			) ([]*models.Picture, error) {
				return implementation.GetAllStickers(ctx)
			},
			ExpectedResult: func() ([]*models.Picture, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request struct{}, m *mocks) {
				m.repository.EXPECT().GetAllStickers(gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*struct{}, error) {
				return &struct{}{}, nil
			},
			Run: func(
				ctx context.Context, implementation *StickerUsecaseImplementation, request struct{},
			) ([]*models.Picture, error) {
				return implementation.GetAllStickers(ctx)
			},
			ExpectedResult: func() ([]*models.Picture, error) {
				return []*models.Picture{&pic}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request struct{}, m *mocks) {
				m.repository.EXPECT().GetAllStickers(gomock.Any()).Return([]*models.Picture{&pic}, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getUseCase(ctrl)
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
			},
		)
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *StickerUsecaseImplementation, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}

type mocks struct {
	repository *MockRepository
}
