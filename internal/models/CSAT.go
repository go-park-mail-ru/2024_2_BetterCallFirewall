package models

type CSAT struct {
	InTotal uint8  `json:"in_total"`
	Review  string `json:"review"`
}

type CSATResult struct {
	InTotalGrade float32 `json:"in_total_grade"`
}
