package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const LIMIT = 10

type CommunityRepository struct {
	db        *sql.DB
	adminList map[uint32][]uint32
}

func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{
		db:        db,
		adminList: make(map[uint32][]uint32),
	}
}

func (c CommunityRepository) GetBatch(ctx context.Context, lastID uint32) ([]*models.CommunityCard, error) {
	var res []*models.CommunityCard
	rows, err := c.db.QueryContext(ctx, GetBatch, lastID, LIMIT)
	if err != nil {
		return nil, fmt.Errorf("get community batch db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		community := &models.CommunityCard{}
		err = rows.Scan(&community.ID, &community.Name, &community.Avatar, &community.About)
		if err != nil {
			return nil, fmt.Errorf("get community rows: %w", err)
		}
		res = append(res, community)
	}
	return res, nil
}

func (c CommunityRepository) GetOne(ctx context.Context, id uint32) (*models.Community, error) {
	res := &models.Community{}
	err := c.db.QueryRowContext(ctx, GetOne, id).Scan(&res.ID, &res.Name, &res.Avatar, &res.About, &res.CountSubscribers)
	if err != nil {
		return nil, fmt.Errorf("get community db: %w", err)
	}
	return res, nil
}

func (c CommunityRepository) Create(ctx context.Context, community *models.Community, author uint32) (uint32, error) {
	res, err := c.db.ExecContext(ctx, CreateNewCommunity, community.Name, community.About, author)
	if err != nil {
		return 0, fmt.Errorf("create community: %w", err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last community id: %w", err)
	}
	id := uint32(lastId)
	c.adminList[id] = append(c.adminList[id], author)
	return id, nil
}

func (c CommunityRepository) Update(ctx context.Context, community *models.Community) error {
	var err error
	if community.Avatar == "" {
		_, err = c.db.ExecContext(ctx, UpdateWithoutAvatar, community.Name, community.About, community.ID)
	} else {
		_, err = c.db.ExecContext(ctx, UpdateWithAvatar, community.Name, community.Avatar, community.About, community.ID)
	}
	if err != nil {
		return fmt.Errorf("update community: %w", err)
	}
	return nil
}

func (c CommunityRepository) Delete(ctx context.Context, id uint32) error {
	_, err := c.db.ExecContext(ctx, Delete, id)
	if err != nil {
		return fmt.Errorf("delete community: %w", err)
	}
	return nil
}

func (c CommunityRepository) JoinCommunity(ctx context.Context, communityId, author uint32) error {
	_, err := c.db.ExecContext(ctx, JoinCommunity, communityId, author)
	if err != nil {
		return fmt.Errorf("join community: %w", err)
	}

	return nil
}

func (c CommunityRepository) LeaveCommunity(ctx context.Context, communityId, author uint32) error {
	_, err := c.db.ExecContext(ctx, LeaveCommunity, communityId, author)
	if err != nil {
		return fmt.Errorf("leave community: %w", err)
	}

	if c.CheckAccess(ctx, communityId, author) {
		admins := c.adminList[communityId]
		var i int
		for idx, admin := range admins {
			if admin == author {
				i = idx
				break
			}
		}
		c.adminList[communityId] = append(c.adminList[communityId][:i], c.adminList[communityId][i+1:]...)
	}

	return nil
}

func (c CommunityRepository) NewAdmin(ctx context.Context, communityId uint32, author uint32) error {
	_, ok := c.adminList[communityId]
	if !ok {
		return my_err.ErrWrongCommunity
	}
	c.adminList[communityId] = append(c.adminList[communityId], author)

	return nil
}

func (c CommunityRepository) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	admins := c.adminList[communityID]
	for _, admin := range admins {
		if admin == userID {
			return true
		}
	}
	return false
}
