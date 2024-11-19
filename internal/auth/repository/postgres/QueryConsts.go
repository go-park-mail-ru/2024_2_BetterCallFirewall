package postgres

const (
	CreateUser     = `INSERT INTO profile (first_name, last_name, email, hashed_password) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING RETURNING id;`
	GetUserByEmail = `SELECT id, first_name, last_name, email, hashed_password FROM profile WHERE email = $1 LIMIT 1;`
)
