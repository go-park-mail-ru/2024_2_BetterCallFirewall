package repository

const (
	InsertNewMetric = `
INSERT INTO csat_metric(total, review, feed) VALUES ($1, $2, $3);`

	GetMetrics = `
SELECT
  (
    SELECT
      COUNT(*)
    FROM
      csat_metric
    WHERE
      total > 3
      AND created_at < cm.created_at AND created_at > $1
  ) * 100 / (
    SELECT
      CASE
        WHEN COUNT(*) = 0 THEN 1
        ELSE COUNT(*)
      END
    FROM
      csat_metric
    WHERE
      created_at < cm.created_at AND created_at > $1
  ) AS grade,
  created_at AS "time"
FROM
  csat_metric AS cm
WHERE
  created_at > $1 AND created_at < $2;`

	GetFeedMetrics = `
SELECT
  (
    SELECT
      COUNT(*)
    FROM
      csat_metric
    WHERE
      feed > 3
      AND created_at < cm.created_at AND created_at > $1
  ) * 100 / (
    SELECT
      CASE
        WHEN COUNT(*) = 0 THEN 1
        ELSE COUNT(*)
      END
    FROM
      csat_metric
    WHERE
      created_at < cm.created_at AND created_at > $1 AND feed <> 0
  ) AS grade,
  created_at AS "time"
FROM
  csat_metric AS cm
WHERE
  created_at > $1 AND created_at < $2;
`
)
