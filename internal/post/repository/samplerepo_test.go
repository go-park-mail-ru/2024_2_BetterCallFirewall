package repository

import (
	"sort"
	"testing"
)

type TestCase struct {
	dataCount int
	lenDate   int
}

func TestRepository(t *testing.T) {
	testCases := []TestCase{
		{dataCount: 10, lenDate: 10},
	}
	for _, testCase := range testCases {
		repo := NewRepository()
		got := repo.GetAll()
		if !sort.SliceIsSorted(got, func(i, j int) bool { return got[i].CreatedAt > got[j].CreatedAt }) {
			t.Errorf("GetAll return not sort by date")
		}
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
