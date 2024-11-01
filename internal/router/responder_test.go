package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TestRouter struct {
	testResponse *httptest.ResponseRecorder
	testData     any
	testErr      error
	expectedCode int
	expectedBody string
	testReqID    string
}

var (
	TestResponder = NewResponder(logrus.New())
	TestError     = errors.New("error")
)

const (
	TestData              = "Sent data to user"
	TestDataOutputBody    = `{"success":true,"data":"Sent data to user"}`
	TestDataWrongMethod   = `{"success":false,"data":"error","message":"method not allowed"}`
	TestDataBadRequest    = `{"success":false,"data":"error","message":"bad request"}`
	TestDataInternalError = `{"success":false,"data":{},"message":"internal server error"}`
	TestDataNoMoreContent = ``
)

func TestOutputJSON(t *testing.T) {
	tests := []TestRouter{
		{
			testResponse: httptest.NewRecorder(),
			testData:     TestData,
			expectedCode: http.StatusOK,
			expectedBody: TestDataOutputBody,
			testReqID:    uuid.New().String(),
		},
	}

	for caseNum, test := range tests {
		TestResponder.OutputJSON(test.testResponse, test.testData, test.testReqID)
		if test.testResponse.Code != test.expectedCode {
			t.Errorf("[%d} wrong status code, expected %d, got %d", caseNum, test.expectedCode, test.testResponse.Code)
		}
		if strings.Compare(test.expectedBody, strings.TrimSpace(test.testResponse.Body.String())) != 0 {
			t.Errorf("[%d] wrong body, expected %s, got %s", caseNum, test.expectedBody, test.testResponse.Body.String())
		}
	}
}

func TestOutputNoMoreContent(t *testing.T) {
	tests := []TestRouter{
		{
			testResponse: httptest.NewRecorder(),
			expectedCode: http.StatusNoContent,
			expectedBody: TestDataNoMoreContent,
			testReqID:    uuid.New().String(),
		},
	}

	for caseNum, test := range tests {
		TestResponder.OutputNoMoreContentJSON(test.testResponse, test.testReqID)
		if test.testResponse.Code != test.expectedCode {
			t.Errorf("[%d} wrong status code, expected %d, got %d", caseNum, test.expectedCode, test.testResponse.Code)
		}
		if strings.Compare(test.expectedBody, strings.TrimSpace(test.testResponse.Body.String())) != 0 {
			t.Errorf("[%d] wrong body, expected %s, got %s", caseNum, test.expectedBody, test.testResponse.Body.String())
		}
	}
}

func TestErrorBadRequest(t *testing.T) {
	tests := []TestRouter{
		{
			testResponse: httptest.NewRecorder(),
			testErr:      TestError,
			expectedCode: http.StatusBadRequest,
			expectedBody: TestDataBadRequest,
			testReqID:    uuid.New().String(),
		},
	}

	for caseNum, test := range tests {
		TestResponder.ErrorBadRequest(test.testResponse, test.testErr, test.testReqID)
		if test.testResponse.Code != test.expectedCode {
			t.Errorf("[%d} wrong status code, expected %d, got %d", caseNum, test.expectedCode, test.testResponse.Code)
		}
		if strings.Compare(test.expectedBody, strings.TrimSpace(test.testResponse.Body.String())) != 0 {
			t.Errorf("[%d] wrong body, expected %s, got %s", caseNum, test.expectedBody, test.testResponse.Body.String())
		}
	}
}

func TestErrorInternal(t *testing.T) {
	tests := []TestRouter{
		{
			testResponse: httptest.NewRecorder(),
			testErr:      TestError,
			expectedCode: http.StatusInternalServerError,
			expectedBody: TestDataInternalError,
			testReqID:    uuid.New().String(),
		},
	}

	for caseNum, test := range tests {
		TestResponder.ErrorInternal(test.testResponse, test.testErr, test.testReqID)
		if test.testResponse.Code != test.expectedCode {
			t.Errorf("[%d} wrong status code, expected %d, got %d", caseNum, test.expectedCode, test.testResponse.Code)
		}
		if strings.Compare(test.expectedBody, strings.TrimSpace(test.testResponse.Body.String())) != 0 {
			t.Errorf("[%d] wrong body, expected %s, got %s", caseNum, test.expectedBody, test.testResponse.Body.String())
		}
	}
}
