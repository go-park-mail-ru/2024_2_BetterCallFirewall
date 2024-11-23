package service

import (
	"context"
	"sync"
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
	isSentCSAT   bool
	spentTime    time.Duration
}

type UERepo struct {
	mapUserExperience map[uint32]*UserExperience
	mutex             sync.RWMutex
}

var (
	repo = UERepo{
		make(map[uint32]*UserExperience),
		sync.RWMutex{},
	}
)

type CSATRepo interface {
	GetMetrics(ctx context.Context, since, before time.Time) (*models.CSATResult, error)
	SaveMetrics(ctx context.Context, csat *models.CSAT) error
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

	if repo.mapUserExperience[userID] == nil {
		repo.mutex.Lock()
		repo.mapUserExperience[userID] = new(UserExperience)
		repo.mutex.Unlock()
		return false
	}

	repo.mutex.RLock()
	isReady := repo.mapUserExperience[userID].addedFriends >= FRIENDS &&
		repo.mapUserExperience[userID].sentMessages >= MESSAGES &&
		repo.mapUserExperience[userID].likes >= LIKES &&
		repo.mapUserExperience[userID].spentTime >= TIME &&
		!repo.mapUserExperience[userID].isSentCSAT
	repo.mutex.RUnlock()
	return isReady
}

func (cs *CSATServiceImpl) SaveMetrics(ctx context.Context, csat *models.CSAT, userID uint32) error {
	err := cs.DB.SaveMetrics(ctx, csat)
	if err != nil {
		return err
	}
	repo.mutex.Lock()
	repo.mapUserExperience[userID].isSentCSAT = true
	repo.mutex.Unlock()
	return nil
}

func (cs *CSATServiceImpl) GetMetrics(ctx context.Context, since, before time.Time) (*models.CSATResult, error) {
	res, err := cs.DB.GetMetrics(ctx, since, before)
	if err != nil {
		return nil, err
	}
	return res, nil
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

func (cs *CSATServiceImpl) NewFriend(userID uint32) {
	repo.mutex.Lock()
	repo.mapUserExperience[userID].NewFriend()
	repo.mutex.Unlock()
}

func (cs *CSATServiceImpl) NewMessage(userID uint32) {
	repo.mutex.Lock()
	repo.mapUserExperience[userID].NewMessage()
	repo.mutex.Unlock()
}

func (cs *CSATServiceImpl) NewLike(userID uint32) {
	repo.mutex.Lock()
	repo.mapUserExperience[userID].NewLike()
	repo.mutex.Unlock()
}

func (cs *CSATServiceImpl) TimeSpent(userID uint32, dur time.Duration) {
	repo.mutex.Lock()
	repo.mapUserExperience[userID].TimeSpent(dur)
	repo.mutex.Unlock()
}
