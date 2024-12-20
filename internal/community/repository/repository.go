package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const LIMIT = 10

type CommunityRepository struct {
	db *sql.DB
}

func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{
		db: db,
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
	err := c.db.QueryRowContext(ctx, GetOne, id).Scan(
		&res.ID, &res.Name, &res.Avatar, &res.About, &res.CountSubscribers,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrWrongCommunity
		}
		return nil, fmt.Errorf("get community db: %w", err)
	}

	return res, nil
}

func (c CommunityRepository) Create(ctx context.Context, community *models.Community, author uint32) (uint32, error) {
	var res *sql.Row

	if community.Avatar == "" {
		res = c.db.QueryRowContext(ctx, CreateNewCommunity, community.Name, community.About, author)
	} else {
		res = c.db.QueryRowContext(
			ctx, CreateNewCommunityWithAvatar, community.Name, community.About, community.Avatar, author,
		)
	}

	err := res.Err()
	if err != nil {
		return 0, fmt.Errorf("create community db: %w", err)
	}

	err = res.Scan(&community.ID)
	if err != nil {
		return 0, fmt.Errorf("create community db: %w", err)
	}

	return community.ID, nil
}

func (c CommunityRepository) Update(ctx context.Context, community *models.Community) error {
	var err error
	if community.Avatar == "" {
		_, err = c.db.ExecContext(ctx, UpdateWithoutAvatar, community.Name, community.About, community.ID)
	} else {
		_, err = c.db.ExecContext(
			ctx, UpdateWithAvatar, community.Name, community.Avatar, community.About, community.ID,
		)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_err.ErrWrongCommunity
		}
		return fmt.Errorf("update community: %w", err)
	}

	return nil
}

func (c CommunityRepository) Delete(ctx context.Context, id uint32) error {
	_, err := c.db.ExecContext(ctx, Delete, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_err.ErrWrongCommunity
		}
		return fmt.Errorf("delete community: %w", err)
	}

	return nil
}

func (c CommunityRepository) JoinCommunity(ctx context.Context, communityId, author uint32) error {
	_, err := c.db.ExecContext(ctx, JoinCommunity, communityId, author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_err.ErrWrongCommunity
		}
		return fmt.Errorf("join community: %w", err)
	}

	return nil
}

func (c CommunityRepository) LeaveCommunity(ctx context.Context, communityId, author uint32) error {
	_, err := c.db.ExecContext(ctx, LeaveCommunity, communityId, author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_err.ErrWrongCommunity
		}
		return fmt.Errorf("leave community: %w", err)
	}
	access := c.CheckAccess(ctx, communityId, author)

	if access {
		_, err := c.db.ExecContext(ctx, DeleteAdmin, communityId, author)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return my_err.ErrWrongCommunity
			}
			return fmt.Errorf("delete admin: %w", err)
		}
	}

	return nil
}

func (c CommunityRepository) NewAdmin(ctx context.Context, communityId uint32, author uint32) error {
	_, err := c.db.ExecContext(ctx, InsertNewAdmin, communityId, author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return my_err.ErrWrongCommunity
		}
		return fmt.Errorf("insert new admin: %w", err)
	}
	return nil
}

func (c CommunityRepository) CheckAccess(ctx context.Context, communityID, userID uint32) bool {
	var num int
	err := c.db.QueryRowContext(ctx, CheckAccess, communityID, userID).Scan(&num)
	if err != nil || num == 0 {
		return false
	}

	return true
}

func (c CommunityRepository) Search(ctx context.Context, query string, lastID uint32) ([]*models.CommunityCard, error) {
	res := make([]*models.CommunityCard, 0)

	rows, err := c.db.QueryContext(ctx, Search, query, lastID, LIMIT)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrNoMoreContent
		}
		return nil, fmt.Errorf("search community: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		community := &models.CommunityCard{}
		err = rows.Scan(&community.ID, &community.Name, &community.Avatar, &community.About)
		if err != nil {
			return nil, fmt.Errorf("search community: %w", err)
		}

		res = append(res, community)
	}

	return res, nil
}

func (c CommunityRepository) GetHeader(ctx context.Context, communityID uint32) (*models.Header, error) {
	row := c.db.QueryRow(GetHeader, communityID)
	header := &models.Header{}
	if row.Err() != nil {
		return nil, my_err.ErrWrongCommunity
	}

	if err := row.Scan(&header.CommunityID, &header.Author, &header.Avatar); err != nil {
		return nil, fmt.Errorf("get header: %w", err)
	}

	return header, nil
}

func (c CommunityRepository) IsFollowed(ctx context.Context, communityID, userID uint32) (bool, error) {
	res := c.db.QueryRowContext(ctx, IsFollow, communityID, userID)
	err := res.Err()
	if err != nil {
		return false, fmt.Errorf("get community db: %w", err)
	}
	var count uint32

	err = res.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("get community db: %w", err)
	}

	return count > 0, nil
}
