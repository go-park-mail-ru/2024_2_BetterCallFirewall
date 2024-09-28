package repository

import (
	"testing"
)

type TestCase struct {
	dataCount int
	lenDate   int
}

func TestRepository(t *testing.T) {
	testCases := []TestCase{
		{dataCount: 0, lenDate: 10},
		{dataCount: 1, lenDate: 10},
		{dataCount: 10, lenDate: 10},
	}
	for _, testCase := range testCases {
		repo := NewRepository()
		repo.FakeData(testCase.dataCount)
		got := repo.GetAll()
		if len(got) != testCase.dataCount {
			t.Errorf("GetAll returned wrong number of results: got %v want %v", len(got), testCase.dataCount)
		}
		for _, val := range got {
			if len(val.CreatedAt) != testCase.lenDate {
				t.Errorf("wrong format Data:%s", val.CreatedAt)
			}
		}
	}
}
