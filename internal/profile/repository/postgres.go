package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

const LIMIT = 20

type ProfileRepo struct {
	DB *sql.DB
}

func NewProfileRepo(db *sql.DB) profile.Repository {
	repo := &ProfileRepo{
		DB: db,
	}
	return repo
}

func (p *ProfileRepo) GetProfileById(ctx context.Context, id uint32) (*models.FullProfile, error) {
	res := &models.FullProfile{}
	err := p.DB.QueryRowContext(ctx, GetProfileByID, id).Scan(&res.ID, &res.FirstName, &res.LastName, &res.Bio, &res.Avatar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErr.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get profile by id db: %w", err)
	}
	return res, nil
}

func (p *ProfileRepo) GetStatus(ctx context.Context, self uint32, profile uint32) (int, error) {
	var status int
	err := p.DB.QueryRowContext(ctx, GetStatus, self, profile).Scan(&status)
	if err != nil {
		return 0, err
	}
	return status, nil
}

func (p *ProfileRepo) GetAll(ctx context.Context, self uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.QueryContext(ctx, GetAllProfilesBatch, self, lastId, LIMIT)
	if err != nil {
		return nil, fmt.Errorf("get all profiles %w", err)
	}
	for rows.Next() {
		profile := &models.ShortProfile{}
		err = rows.Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get all profiles db: %w", err)
		}
		res = append(res, profile)
	}
	rows.Close()
	return res, nil
}

func (p *ProfileRepo) UpdateProfile(profile *models.FullProfile) error {
	_, err := p.DB.Exec(UpdateProfile, profile.FirstName, profile.LastName, profile.Bio, profile.ID)
	if err != nil {
		return fmt.Errorf("update profile %w", err)
	}

	return nil
}

func (p *ProfileRepo) DeleteProfile(u uint32) error {
	_, err := p.DB.Exec(DeleteProfile, u)
	if err != nil {
		return fmt.Errorf("delete profile %w", err)
	}
	return nil
}

func (p *ProfileRepo) AddFriendsReq(receiver uint32, sender uint32) error {
	_, err := p.DB.Exec(AddFriends, sender, receiver)
	if err != nil {
		return fmt.Errorf("add friend db: %w", err)
	}
	return nil
}

func (p *ProfileRepo) AcceptFriendsReq(who uint32, whose uint32) error {
	_, err := p.DB.Exec(AcceptFriendReq, whose, who)
	if err != nil {
		return fmt.Errorf("accept friend db: %w", err)
	}
	return nil
}

func (p *ProfileRepo) MoveToSubs(who uint32, whom uint32) error {
	_, err := p.DB.Exec(RemoveFriendsReq, who, whom)
	if err != nil {
		return fmt.Errorf("remove friends db: %w", err)
	}
	return nil
}

func (p *ProfileRepo) RemoveSub(who uint32, whom uint32) error {
	_, err := p.DB.Exec(DeleteFriendship, who, whom)
	if err != nil {
		return fmt.Errorf("delete friendship db: %w", err)
	}
	return nil
}

func (p *ProfileRepo) GetAllFriends(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.QueryContext(ctx, GetAllFriends, u, lastId, LIMIT)
	if err != nil {
		return nil, fmt.Errorf("get all friends %w", err)
	}
	for rows.Next() {
		profile := models.ShortProfile{}
		err = rows.Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get all friends db: %w", err)
		}
		res = append(res, &profile)
	}
	return res, nil
}

func (p *ProfileRepo) GetAllSubs(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.QueryContext(ctx, GetAllSubs, u, lastId, LIMIT)
	if err != nil {
		return nil, fmt.Errorf("get all subs db: %w", err)
	}
	for rows.Next() {
		profile := models.ShortProfile{}
		err = rows.Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get all subs db: %w", err)
		}
		res = append(res, &profile)
	}
	return res, nil
}

func (p *ProfileRepo) GetAllSubscriptions(ctx context.Context, u uint32, lastId uint32) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.QueryContext(ctx, GetAllSubscriptions, u, lastId, LIMIT)
	if err != nil {
		return nil, fmt.Errorf("get all subscriptions db: %w", err)
	}
	for rows.Next() {
		profile := models.ShortProfile{}
		err = rows.Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get all subscriptions db: %w", err)
		}
		res = append(res, &profile)
	}
	return res, nil
}

func (p *ProfileRepo) GetFriendsID(ctx context.Context, u uint32) ([]uint32, error) {
	res := make([]uint32, 0)
	rows, err := p.DB.QueryContext(ctx, GetFriendsID, u)
	if err != nil {
		return nil, fmt.Errorf("get friends id db: %w", err)
	}
	for rows.Next() {
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("get friends id db: %w", err)
		}
		res = append(res, id)
	}
	return res, nil
}

func (p *ProfileRepo) GetSubscribersID(ctx context.Context, u uint32) ([]uint32, error) {
	res := make([]uint32, 0)
	rows, err := p.DB.QueryContext(ctx, GetSubsID, u)
	if err != nil {
		return nil, fmt.Errorf("get friends id db: %w", err)
	}
	for rows.Next() {
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("get friends id db: %w", err)
		}
		res = append(res, id)
	}
	return res, nil
}

func (p *ProfileRepo) GetSubscriptionsID(ctx context.Context, u uint32) ([]uint32, error) {
	res := make([]uint32, 0)
	rows, err := p.DB.QueryContext(ctx, GetSubscriptionsID, u)
	if err != nil {
		return nil, fmt.Errorf("get friends id db: %w", err)
	}
	for rows.Next() {
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("get friends id db: %w", err)
		}
		res = append(res, id)
	}
	return res, nil
}

func (p *ProfileRepo) GetHeader(ctx context.Context, u uint32) (*models.Header, error) {
	profile := &models.Header{AuthorID: u}
	err := p.DB.QueryRowContext(ctx, GetShortProfile, u).Scan(&profile.Author, &profile.Avatar)
	if err != nil {
		return nil, fmt.Errorf("get short profile db: %w", err)
	}
	return profile, nil
}
