package postgres

const (
	CreateUser            = `INSERT INTO profile (first_name, last_name, email, hashed_password) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING RETURNING id;`
	GetUserByEmail        = `SELECT id, first_name, last_name, email, hashed_password FROM profile WHERE email = $1 LIMIT 1;`
	CreateNewSessionTable = `CREATE TABLE IF NOT EXISTS session (id SERIAL PRIMARY KEY, sess_id TEXT NOT NULL, user_id INTEGER NOT NULL UNIQUE, created_at BIGINT NOT NULL);`
	CreateSession         = `INSERT INTO session (sess_id, user_id, created_at) VALUES ($1, $2, $3) ON CONFLICT(user_id) DO UPDATE SET sess_id = EXCLUDED.sess_id, created_at = EXCLUDED.created_at;`
	FindSession           = `SELECT sess_id, user_id, created_at FROM session WHERE sess_id = $1 LIMIT 1;`
	DeleteSession         = `DELETE FROM session WHERE sess_id = $1;`
	DeleteOutdatedSession = `DELETE FROM session WHERE created_at <= $1;`
)
