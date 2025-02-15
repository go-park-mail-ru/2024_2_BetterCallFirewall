package repository

const (
	InsertNewSticker = `INSERT INTO sticker(file_path, profile_id) VALUES ($1, $2)`
	GetAllSticker    = `SELECT file_path FROM sticker`
	GetUserStickers  = `SELECT file_path FROM sticker WHERE user_id = $1`
)
