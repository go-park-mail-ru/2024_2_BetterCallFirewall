package postgres

const (
	getAllChatBatch = `WITH last_messages AS (
    SELECT
        CASE
            WHEN sender = $1 THEN receiver
            ELSE sender
        END AS related_user,
        content,
        created_at,
        ROW_NUMBER() OVER (PARTITION BY
            CASE
                WHEN sender = $1 THEN receiver
                ELSE sender
                END
            ORDER BY created_at DESC) AS rn
    FROM
        message
    WHERE
        sender = $1 OR receiver = $1
)

SELECT
   	related_user,
    profile.first_name || ' ' || profile.last_name AS chat,
    content AS last_message_content,
    last_message.created_at AS last_message_time
FROM
    last_messages INNER JOIN profile ON related_user = profile.id
WHERE
    rn = 1 AND last_message_time < $2
ORDER BY
    last_message_time DESC
LIMIT 15;`

	getLatestMessagesBatch = `SELECT sender, content, created_at
FROM message
WHERE ((sender = $1 AND receiver = $2) OR (sender = $2 AND receiver = $1)) 
AND created_at < $3
ORDER BY created_at DESC
LIMIT 20;`

	sendNewMessage = `INSERT INTO message(receiver, sender, content) VALUES ($1, $2, $3)`
)
