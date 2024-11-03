package repository

const (
	InsertPostFile = `INSERT INTO file(file_path, profile_id , post_id) VALUES ($1, $2, $3);`
	GetPostFile    = `SELECT file_path FROM file WHERE post_id = $1 LIMIT 1;`
	GetProfileFile = `SELECT file_path FROM file WHERE profile_id = $1;`
)
