package repository

const (
	AddLikeToPost         = `INSERT INTO reaction (post_id, user_id) VALUES ($1, $2);`
	AddLikeToComment      = `INSERT INTO reaction (comment_id, user_id) VALUES ($1, $2);`
	AddLikeToFile         = `INSERT INTO reaction (file_id, user_id) VALUES ($1, $2);`
	DeleteLikeFromPost    = `DELETE FROM reaction WHERE post_id = $1 AND user_id = $2;`
	DeleteLikeFromComment = `DELETE FROM reaction WHERE comment_id = $1 AND user_id = $2;`
	DeleteLikeFromFile    = `DELETE FROM reaction WHERE file_id = $1 AND user_id = $2;`
	GetLikesOnPost        = `SELECT COUNT(*) FROM reaction WHERE post_id = $1;`
)
