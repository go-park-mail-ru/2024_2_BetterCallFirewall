package service

import (
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

const (
	FRIENDS  = 10
	MESSAGES = 10
	TIME     = time.Second * 24
	LIKES    = 10
)

type UserExperience struct {
	addedFriends uint32
	sentMessages uint32
	likes        uint32
	spentTime    time.Duration
}

func (ue *UserExperience) NewFriend() {
	ue.addedFriends++
}

func (ue *UserExperience) NewMessage() {
	ue.sentMessages++
}

func (ue *UserExperience) NewLike() {
	ue.likes++
}

func (ue *UserExperience) TimeSpent(sessTiming time.Duration) {
	ue.spentTime += sessTiming
}

var mapUserExperiences = make(map[uint32]*UserExperience)

type CSATRepo interface {
	GetMetrics(since, before time.Time) (float32, error)
	SaveMetrics(csat models.CSAT) error
}

type CSATServiceImpl struct {
	DB CSATRepo
}

func NewCSATServiceImpl(db CSATRepo) *CSATServiceImpl {
	return &CSATServiceImpl{
		DB: db,
	}
}

func (cs *CSATServiceImpl) CheckExperience(userID uint32) bool {
	if mapUserExperiences[userID] == nil {
		mapUserExperiences[userID] = new(UserExperience)
		return false
	}
	isReady := mapUserExperiences[userID].addedFriends >= FRIENDS &&
		mapUserExperiences[userID].sentMessages >= MESSAGES &&
		mapUserExperiences[userID].likes >= LIKES &&
		mapUserExperiences[userID].spentTime >= TIME

	return isReady
}

func (cs *CSATServiceImpl) SaveMetrics(csat *models.CSAT) error {
	err := cs.SaveMetrics(csat)
	if err != nil {
		return err
	}
	return nil
}
