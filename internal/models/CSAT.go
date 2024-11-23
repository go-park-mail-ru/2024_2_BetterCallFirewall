package models

type CSAT struct {
	InTotal uint8  `json:"in_total"`
	Review  string `json:"review"`
	Feed    uint8  `json:"feed"`
}

type CSATResult struct {
	InTotalGrade float32 `json:"in_total_grade"`
	FeedGrade    float32 `json:"feed_grade"`
}
