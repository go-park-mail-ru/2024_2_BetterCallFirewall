package service

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

var errMock = errors.New("something with DB")

type MockDB struct{}

func (m MockDB) Create(user *models.User) (uint32, error) {
	if user.ID == 0 {
		return user.ID, myErr.ErrUserNotFound
	}
	return user.ID, nil
}

func (m MockDB) GetByEmail(email string) (*models.User, error) {
	if email == "email@wrong.com" {
		return nil, myErr.ErrUserNotFound
	}

	if email == "email@wrong2.com" {
		return nil, errMock
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	return &models.User{
		Email:    email,
		Password: string(hash),
	}, nil
}

type TestCase struct {
	user      models.User
	wantError error
}

func TestCreate(t *testing.T) {
	serv := NewAuthServiceImpl(MockDB{})

	testCases := []TestCase{
		{models.User{ID: 1, Email: "email@email.com", Password: "some password"}, nil},
		{models.User{ID: 0, Email: "email@email.com", Password: "some password"}, myErr.ErrUserNotFound},
		{models.User{ID: 100, Email: "email", Password: "some password"}, myErr.ErrNonValidEmail},
		{models.User{ID: 1, Email: "email@email.com", Password: "some very very long password, more then 74 symbols this password dont use anymore in real life and have validate on client"},
			bcrypt.ErrPasswordTooLong},
	}

	for _, testCase := range testCases {
		_, err := serv.Register(testCase.user)
		if !errors.Is(err, testCase.wantError) {
			t.Errorf("Register() error = %v, wantErr %v", err, testCase.wantError)
		}
	}
}

func TestAuth(t *testing.T) {
	serv := NewAuthServiceImpl(MockDB{})

	testCases := []TestCase{
		{models.User{ID: 1, Email: "email@email.com", Password: "password"}, nil},
		{models.User{ID: 100, Email: "email", Password: "some password"}, myErr.ErrNonValidEmail},
		{models.User{ID: 100, Email: "email@wrong.com", Password: "some password"}, myErr.ErrWrongEmailOrPassword},
		{models.User{ID: 100, Email: "email@wrong2.com", Password: "some password"}, errMock},
		{models.User{ID: 1, Email: "email@email.com", Password: "some password"}, myErr.ErrWrongEmailOrPassword},
	}

	for _, testCase := range testCases {
		_, err := serv.Auth(testCase.user)
		if !errors.Is(err, testCase.wantError) {
			t.Errorf("Auth() error = %v, wantErr %v", err, testCase.wantError)
		}
	}
}

type TestCaseValidate struct {
	email string
	pass  bool
}

func TestValidateEmail(t *testing.T) {
	serv := NewAuthServiceImpl(MockDB{})

	testCases := []TestCaseValidate{
		{email: "email@email.com", pass: true},
		{email: "loop@mail.ru", pass: true},
		{email: "loop", pass: false},
		{email: "", pass: false},
		{email: "@email.com", pass: false},
		{email: "loop@mailru", pass: false},
		{email: "email-my@mail.ru", pass: true},
	}

	for _, testCase := range testCases {
		res := serv.validateEmail(testCase.email)
		if res != testCase.pass {
			t.Errorf("ValidateEmail() error = %v, wantErr %v, with email: %v", res, testCase.pass, testCase.email)
		}
	}
}
