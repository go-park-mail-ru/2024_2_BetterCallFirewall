package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type ProfileRepo struct {
	DB *sql.DB
}

func NewProfileRepo(db *sql.DB) profile.Repository {
	repo := &ProfileRepo{
		DB: db,
	}
	return repo
}

func (p *ProfileRepo) GetProfileById(id uint32, ctx context.Context) (*models.FullProfile, error) {
	res := &models.FullProfile{}
	err := p.DB.QueryRowContext(ctx, GetProfileByID, id).Scan(&res.ID, &res.FirstName, &res.LastName, &res.Bio, &res.Avatar)
	if err != nil {
		return nil, fmt.Errorf("get profile by id db: %w", err)
	}
	return res, nil
}

func (p *ProfileRepo) GetAll(self uint32, ctx context.Context) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.QueryContext(ctx, GetAllProfiles, self)
	if err != nil {
		return nil, fmt.Errorf("get all profiles %w", err)
	}
	for rows.Next() {
		profile := &models.ShortProfile{}
		err = rows.Scan(profile.ID, profile.FirstName, profile.LastName, profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get all profiles db: %w", err)
		}
		res = append(res, profile)
	}
	rows.Close()
	return res, nil
}

func (p *ProfileRepo) UpdateProfile(profile *models.FullProfile) error {
	_, err := p.DB.Exec(UpdateProfile, profile.FirstName, profile.LastName, profile.Bio, profile.Avatar, profile.ID)
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

func (p *ProfileRepo) GetAllFriends(u uint32, ctx context.Context) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.QueryContext(ctx, GetAllFriends, u)
	if err != nil {
		return nil, fmt.Errorf("get all friends %w", err)
	}
	for rows.Next() {
		profile := &models.ShortProfile{}
		err = rows.Scan(profile.ID, profile.FirstName, profile.LastName, profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get all friends db: %w", err)
		}
		res = append(res, profile)
	}
	return res, nil
}

func (p *ProfileRepo) CheckStatus(u1 uint32, u2 uint32, ctx context.Context) (int, error) {
	var status int
	err := p.DB.QueryRowContext(ctx, CheckFriendReq, u1, u2).Scan(&status)
	if err != nil {
		return status, fmt.Errorf("check status db: %w", err)
	}
	return status, nil
}

func (p *ProfileRepo) GetFriendsID(u uint32, ctx context.Context) ([]uint32, error) {
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

func (p *ProfileRepo) GetHeader(u uint32) (*models.Header, error) {
	profile := &models.Header{AuthorID: u}
	err := p.DB.QueryRowContext(context.Background(), GetShortProfile, u).Scan(&profile.Author, &profile.Avatar)
	if err != nil {
		return nil, fmt.Errorf("get short profile db: %w", err)
	}
	return profile, nil
}
