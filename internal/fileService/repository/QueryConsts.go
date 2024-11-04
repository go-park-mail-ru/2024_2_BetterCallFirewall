package repository

const (
	InsertPostFile    = `INSERT INTO file(file_path, post_id) VALUES ($1, $2);`
	InsertProfileFile = `INSERT INTO file(file_path, profile_id) VALUES ($1, $2);`
	GetPostFile       = `SELECT file_path FROM file WHERE post_id = $1 LIMIT 1;`
	GetProfileFile    = `SELECT file_path FROM file WHERE profile_id = $1;`
	UpdatePostFile    = `UPDATE file SET file_path = $1 WHERE post_id = $2;`
)
