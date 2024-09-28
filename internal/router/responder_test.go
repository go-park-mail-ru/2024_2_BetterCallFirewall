package router

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type TestRouter struct {
	testResponse *httptest.ResponseRecorder
	testData     any
	testErr      error
	expectedCode int
	expectedBody string
}

var (
	TestResponder = NewResponder(log.New(os.Stdout, "", log.LstdFlags))
	TestError     = errors.New("error")
)

const (
	TestData              = "Sent data to user"
	TestDataOutputBody    = `{"success":true,"data":"Sent data to user","message":""}`
	TestDataWrongMethod   = `{"success":false,"data":"error","message":"method not allowed"}`
	TestDataBadRequest    = `{"success":false,"data":"error","message":"bad request"}`
	TestDataInternalError = `{"success":false,"data":"error","message":"internal server error"}`
)

func TestOutputJSON(t *testing.T) {
	tests := []TestRouter{
		{
			testResponse: httptest.NewRecorder(),
			testData:     TestData,
			expectedCode: http.StatusOK,
			expectedBody: TestDataOutputBody,
		},
	}

	for caseNum, test := range tests {
		TestResponder.OutputJSON(test.testResponse, test.testData)
		if test.testResponse.Code != test.expectedCode {
			t.Errorf("[%d} wrong status code, expected %d, got %d", caseNum, test.expectedCode, test.testResponse.Code)
		}
		if strings.Compare(test.expectedBody, strings.TrimSpace(test.testResponse.Body.String())) != 0 {
			t.Errorf("[%d] wrong body, expected %s, got %s", caseNum, test.expectedBody, test.testResponse.Body.String())
		}
	}
}

func TestErrorWrongMethod(t *testing.T) {
	tests := []TestRouter{
		{
			testResponse: httptest.NewRecorder(),
			testErr:      TestError,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: TestDataWrongMethod,
		},
	}

	for caseNum, test := range tests {
		TestResponder.ErrorWrongMethod(test.testResponse, test.testErr)
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
		},
	}

	for caseNum, test := range tests {
		TestResponder.ErrorBadRequest(test.testResponse, test.testErr)
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
		},
	}

	for caseNum, test := range tests {
		TestResponder.ErrorInternal(test.testResponse, test.testErr)
		if test.testResponse.Code != test.expectedCode {
			t.Errorf("[%d} wrong status code, expected %d, got %d", caseNum, test.expectedCode, test.testResponse.Code)
		}
		if strings.Compare(test.expectedBody, strings.TrimSpace(test.testResponse.Body.String())) != 0 {
			t.Errorf("[%d] wrong body, expected %s, got %s", caseNum, test.expectedBody, test.testResponse.Body.String())
		}
	}
}
