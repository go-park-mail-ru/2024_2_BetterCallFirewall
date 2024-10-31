package postgres

const (
	getAllChatBatch = `SELECT sender, content, created_at 
FROM message
WHERE (receiver, created_at) IN
    (
    SELECT receiver, MAX(created_at) AS last_time
    FROM message
    WHERE created_at < $1 AND (receiver=$2 or sender=$2)
    GROUP BY receiver
	)
ORDER BY created_at DESC 
LIMIT 15;`
)
