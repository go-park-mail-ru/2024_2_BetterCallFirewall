package repository

const (
	CreateUser     = `INSERT INTO profile (first_name, last_name, email, hashed_password) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING RETURNING id;`
	GetUserByEmail = `SELECT id, first_name, last_name, email, hashed_password FROM profile WHERE email = $1 LIMIT 1;`

	GetProfileByID      = "SELECT profile.id, first_name, last_name, bio, avatar FROM profile WHERE profile.id = $1 LIMIT 1;"
	GetStatus           = "SELECT status FROM friend WHERE (sender = $1 AND receiver = $2) LIMIT 1"
	GetAllProfilesBatch = "WITH friends AS (SELECT sender AS friend FROM friend WHERE (receiver = $1 AND status = 0) UNION SELECT receiver AS friend FROM friend WHERE (sender = $1 AND status = 0)), subscriptions AS (SELECT sender AS subscription FROM friend WHERE (receiver = $1 AND status = -1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = 1)) SELECT p.id, first_name, last_name, avatar FROM profile p WHERE p.id <> $1 AND p.id > $2 AND p.id NOT IN (SELECT friend FROM friends) AND p.id NOT IN (SELECT subscription FROM subscriptions) ORDER BY p.id LIMIT $3;"
	UpdateProfile       = "UPDATE profile SET first_name = $1, last_name = $2, bio = $3 WHERE id = $4;"
	UpdateProfileAvatar = "UPDATE profile SET avatar = $2, first_name = $3, last_name = $4, bio = $5 WHERE id = $1;"
	DeleteProfile       = "DELETE FROM profile WHERE id = $1;"
	AddFriends          = "INSERT INTO friend(sender, receiver, status) VALUES ($1, $2, 1);"
	AcceptFriendReq     = "UPDATE friend SET status = 0 WHERE sender = $1 AND receiver = $2;"
	RemoveFriendsReq    = "UPDATE friend SET status = ( CASE WHEN sender = $1 THEN -1 ELSE 1 END) WHERE (receiver = $1 AND sender = $2) OR (sender = $1 AND receiver = $2);"
	GetAllFriends       = "WITH friends AS (SELECT sender AS friend FROM friend WHERE (receiver = $1 AND status = 0) UNION SELECT receiver AS friend FROM friend WHERE (sender = $1 AND status = 0)) SELECT profile.id, first_name, last_name, avatar FROM profile INNER JOIN friends ON friend = profile.id WHERE profile.id > $2 ORDER BY profile.id LIMIT $3;"
	GetAllSubs          = "WITH subs AS ( SELECT sender AS subscriber FROM friend WHERE (receiver = $1 AND status = 1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = -1)) SELECT profile.id, first_name, last_name, avatar FROM profile INNER JOIN subs ON subscriber = profile.id WHERE profile.id > $2 ORDER BY profile.id LIMIT $3;"
	GetAllSubscriptions = "WITH subscriptions AS ( SELECT sender AS subscription FROM friend WHERE (receiver = $1 AND status = -1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = 1)) SELECT profile.id, first_name, last_name, avatar FROM profile INNER JOIN subscriptions ON subscription = profile.id WHERE profile.id > $2 ORDER BY profile.id LIMIT $3;"
	CheckFriendship     = `SELECT status FROM friend WHERE sender = $2 AND receiver = $1;`

	DeleteFriendship = "DELETE FROM friend WHERE (sender = $1 AND receiver = $2) OR (receiver = $1 AND sender = $2);"

	GetFriendsID       = "SELECT sender AS friend FROM friend WHERE (receiver = $1 AND status = 0) UNION SELECT receiver AS friend FROM friend WHERE (sender = $1 AND status = 0)"
	GetSubsID          = "SELECT sender AS subscriber FROM friend WHERE (receiver = $1 AND status = 1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = -1)"
	GetSubscriptionsID = "SELECT sender AS subscription FROM friend WHERE (receiver = $1 AND status = -1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = 1)"
	GetAllStatuses     = "WITH friends AS (\n    SELECT sender AS friend\n    FROM friend\n    WHERE (receiver = $1 AND status = 0)\n    UNION\n    SELECT receiver AS friend\n    FROM friend\n    WHERE (sender = $1 AND status = 0)\n), subscriptions AS (\n    SELECT sender AS subscription FROM friend WHERE (receiver = $1 AND status = -1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = 1)\n), subscribers AS (\n    SELECT sender AS subscriber FROM friend WHERE (receiver = $1 AND status = 1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = -1)) SELECT (SELECT json_agg(friend) FROM friends) AS friends, (SELECT json_agg(subscriber) FROM subscribers) AS subscribers, (SELECT json_agg(subscription) FROM subscriptions) AS subscriptions;"
	GetShortProfile    = "SELECT first_name || ' ' || last_name AS name, avatar FROM profile WHERE profile.id = $1 LIMIT 1;"

	GetCommunitySubs = `WITH subs AS (SELECT profile_id AS id FROM community_profile WHERE community_id = $1) SELECT p.id, first_name, last_name, avatar FROM profile p JOIN subs ON p.id = subs.id WHERE id > $2 ORDER BY id LIMIT $3;`
)
