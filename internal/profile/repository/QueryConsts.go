package repository

const (
	GetProfileByID      = "SELECT id, first_name, last_name, bio, avatar FROM profile WHERE id = $1 LIMIT 1;"
	GetAllProfiles      = "SELECT id, first_name, last_name, avatar FROM profile WHERE id <> $1 ;"
	UpdateProfile       = "UPDATE profile SET first_name = $1, last_name = $2, bio = $3, avatar = $4 WHERE id = $5;"
	DeleteProfile       = "DELETE FROM profile WHERE id = $1;"
	AddFriends          = "INSERT INTO friend(sender, receiver, status) VALUES ($1, $2, 1);"
	AcceptFriendReq     = "UPDATE friend SET status = 0 WHERE sender = $1 AND receiver = $2;"
	RemoveFriendsReq    = "UPDATE friend SET status = ( CASE WHEN sender = $1 THEN -1 ELSE 1 END) WHERE (receiver = $1 AND sender = $2) OR (sender = $1 AND receiver = $2);"
	GetAllFriends       = "WITH friendships AS ( SELECT sender, receiver FROM friend WHERE (sender = $1 OR receiver = $1) AND status = 0), friends AS ( SELECT CASE WHEN sender = $1 THEN receiver ELSE sender END AS friend_id FROM friendships) SELECT profile.id, first_name, last_name FROM friends INNER JOIN profile ON friend_id = profile.id LEFT JOIN file ON profile.avatar = file.id;"
	GetAllSubs          = "WITH friendships AS ( SELECT sender, receiver FROM friend WHERE (receiver = $1 AND status = 1) OR (receiver = $1 AND status = -1), friends AS ( SELECT CASE WHEN sender = $1 THEN receiver ELSE sender END AS friend_id FROM friendships) SELECT profile.id, first_name, last_name FROM friends INNER JOIN profile ON friend_id = profile.id LEFT JOIN file ON profile.avatar = file.id;"
	GetAllSubscriptions = "WITH friendships AS ( SELECT sender, receiver FROM friend WHERE (receiver = $1 AND status = -1) OR (receiver = $1 AND status = 1), friends AS ( SELECT CASE WHEN sender = $1 THEN receiver ELSE sender END AS friend_id FROM friendships) SELECT profile.id, first_name, last_name FROM friends INNER JOIN profile ON friend_id = profile.id LEFT JOIN file ON profile.avatar = file.id;"

	DeleteFriendship = "DELETE FROM friend WHERE (sender = $1 AND receiver = $2) OR (receiver = $1 AND sender = $2);"

	GetFriendsID    = "WITH friendships AS ( SELECT sender, receiver FROM friend WHERE (sender = $1 OR receiver = $1) AND status = 0) SELECT CASE WHEN sender = $1 THEN receiver ELSE sender END AS friend_id FROM friendships;"
	GetShortProfile = "SELECT first_name || ' ' || last_name AS name, avatar FROM profile WHERE id = $1 LIMIT 1;"
)
