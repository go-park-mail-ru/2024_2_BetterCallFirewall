package start_postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sirupsen/logrus"
)

var connect = map[string]*sql.DB{"host=test port=test user=test password=test dbname=test sslmode=test": {}}

func StartPostgres(connStr string, logger *logrus.Logger) (*sql.DB, error) {
	if _, ok := connect[connStr]; ok {
		return connect[connStr], nil
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	db.SetMaxOpenConns(10)

	retrying := 10
	i := 1
	logger.Infof("try ping postgresql:%v", i)
	for err = db.Ping(); err != nil; err = db.Ping() {
		if i >= retrying {
			return nil, fmt.Errorf("postgres connect: %w", err)
		}
		i++
		time.Sleep(1 * time.Second)
		logger.Infof("try ping postgresql: %v", i)
	}

	connect[connStr] = db
	return db, nil
}
