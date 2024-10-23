package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type SessionRedisRepository struct {
	db redis.Conn
}

func NewSessionRedisRepository(db redis.Conn) *SessionRedisRepository {
	return &SessionRedisRepository{
		db: db,
	}
}

func (s *SessionRedisRepository) CreateSession(session *models.Session) error {
	dataSerialized, err := json.Marshal(session)
	if err != nil {
		return err
	}
	mkey := "sessions:" + session.ID
	res, err := redis.String(s.db.Do("SET", mkey, dataSerialized, "EX", 86400))
	if err != nil {
		return err
	}
	if res != "OK" {
		return myErr.ErrResNotOK
	}
	return nil
}

func (s *SessionRedisRepository) FindSession(sessID string) (*models.Session, error) {
	mkey := "sessions:" + sessID
	data, err := redis.Bytes(s.db.Do("GET", mkey))
	if err != nil {
		return nil, myErr.ErrSessionNotFound
	}
	sess := &models.Session{}
	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, err
	}
	err = s.UpdateSession(sessID)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *SessionRedisRepository) DestroySession(sessID string) error {
	mkey := "sessions:" + sessID
	_, err := redis.Int(s.db.Do("DEL", mkey))
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionRedisRepository) UpdateSession(sessID string) error {
	mkey := "sessions:" + sessID
	_, err := redis.String(s.db.Do("EXPIRE", mkey, 86400))
	if err != nil {
		return err
	}
	return nil
}
