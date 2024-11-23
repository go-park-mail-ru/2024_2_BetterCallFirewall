package repository

const (
	InsertNewMetric = `
INSERT INTO csat_metric(total, review) VALUES ($1, $2);`

	GetMetrics = `
WITH good_total_reviews AS (
	SELECT COUNT(*) 
	FROM csat_metric
	WHERE created_at > $1 AND created_at < $2
	GROUP BY total 
	HAVING total > 3
) SELECT COUNT(*) * 100 / (SELECT * FROM good_total_reviews) AS total_grade 
FROM csat_metric
WHERE created_at > $1 AND created_at < $2;
`
)
