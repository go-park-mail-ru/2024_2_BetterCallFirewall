package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type SessionRedisRepository struct {
	db *redis.Pool
}

func NewSessionRedisRepository(db *redis.Pool) *SessionRedisRepository {
	return &SessionRedisRepository{
		db: db,
	}
}

func (s *SessionRedisRepository) CreateSession(session *models.Session) error {
	conn := s.db.Get()
	defer conn.Close()
	dataSerialized, err := json.Marshal(session)
	if err != nil {
		return err
	}
	mkey := "sessions:" + session.ID

	res, err := redis.String(conn.Do("SET", mkey, dataSerialized, "EX", 86400))
	if err != nil {
		return err
	}

	if res != "OK" {
		return myErr.ErrResNotOK
	}

	return nil
}

func (s *SessionRedisRepository) FindSession(sessID string) (*models.Session, error) {
	conn := s.db.Get()
	defer conn.Close()
	mkey := "sessions:" + sessID
	data, err := redis.String(conn.Do("GET", mkey))
	if err != nil {
		return nil, myErr.ErrSessionNotFound
	}

	sess := &models.Session{}
	err = json.Unmarshal([]byte(data), sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *SessionRedisRepository) DestroySession(sessID string) error {
	conn := s.db.Get()
	defer conn.Close()
	mkey := "sessions:" + sessID
	_, err := redis.Int(conn.Do("DEL", mkey))
	if err != nil {
		return err
	}
	return nil
}
