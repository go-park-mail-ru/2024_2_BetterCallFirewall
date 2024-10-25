package repository

const (
	GetProfileByID      = "SELECT id, first_name, last_name, bio, avatar FROM profile LEFT JOIN file ON profile.avatar = file.id WHERE id = $1 LIMIT 1;"
	GetAllProfiles      = "SELECT id, first_name, last_name, avatar FROM profile LEFT JOIN file ON profile.avatar = file.id WHERE id <> $1 ;"
	UpdateProfile       = "UPDATE profile SET first_name = $1, last_name = $2, bio = $3 WHERE id = $4;"
	UpdateProfileAvatar = "WITH new_avatar AS (INSERT INTO file(profile_id, file_path) VALUES ($1, $2) RETURNING id) UPDATE profile SET avatar = (SELECT id FROM new_avatar) WHERE id = $1;"
	DeleteProfile       = "DELETE FROM profile WHERE id = $1;"
	AddFriends          = "INSERT INTO friend(sender, receiver, status) VALUES ($1, $2, 1);"
	AcceptFriendReq     = "UPDATE friend SET status = 0 WHERE sender = $1 AND receiver = $2;"
	RemoveFriendsReq    = "UPDATE friend SET status = ( CASE WHEN sender = $1 THEN -1 ELSE 1 END) WHERE (receiver = $1 AND sender = $2) OR (sender = $1 AND receiver = $2);"
	GetAllFriends       = "WITH friends AS (SELECT sender AS friend FROM friend WHERE (receiver = $1 AND status = 0) UNION SELECT receiver AS friend FROM friend WHERE (sender = $1 AND status = 0)) SELECT id, first_name, last_name FROM profile INNER JOIN friends ON friend = profile.id LEFT JOIN file ON profile.avatar = file.id;"
	GetAllSubs          = "WITH subs AS ( SELECT sender AS subscriber FROM friend WHERE (receiver = $1 AND status = 1) UNION SELECT receiver AS subscriber FROM friend WHERE (sender = $1 AND status = -1)) SELECT profile_id, first_name, last_name FROM profile INNER JOIN subs ON subscriber = profile.id LEFT JOIN file ON profile.avatar = file.id;"
	GetAllSubscriptions = "WITH subscriptions AS ( SELECT sender AS subscription FROM friend WHERE (receiver = $1 AND status = -1) UNION SELECT receiver AS subscriber FROM subscription WHERE (sender = $1 AND status = 1)) SELECT profile_id, first_name, last_name FROM profile INNER JOIN subscriptions ON subscription = profile.id LEFT JOIN file ON profile.avatar = file.id;"

	DeleteFriendship = "DELETE FROM friend WHERE (sender = $1 AND receiver = $2) OR (receiver = $1 AND sender = $2);"

	GetFriendsID    = "SELECT sender AS friend FROM friend WHERE (receiver = $1 AND status = 0) UNION SELECT receiver AS friend FROM friend WHERE (sender = $1 AND status = 0)"
	GetShortProfile = "SELECT first_name || ' ' || last_name AS name, avatar FROM profile WHERE id = $1 LIMIT 1;"
)
