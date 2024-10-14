package postgres

const (
	GetProfileByID = "SELECT id, first_name, last_name, bio, avatar FROM profile WHERE id = $1;"
	GetAllProfiles = "SELECT id, first_name, last_name, avatar FROM person WHERE id != $1 ;"
	UpdateProfile  = "UPDATE profile SET first_name = $1, last_name = $2, bio = $3, avatar = $4 WHERE id = $5;"
	DeleteProfile  = "DELETE FROM profile WHERE id = $1;"
	AddFriends     = ""
)
