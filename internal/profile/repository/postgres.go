package repository

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type ProfileRepo struct {
	DB *pgx.ConnPool
}

func NewProfileRepo(db *pgx.ConnPool) profile.Repository {
	repo := &ProfileRepo{
		DB: db,
	}
	return repo
}

func (p *ProfileRepo) GetProfileById(id uint32) (*models.FullProfile, error) {
	res := &models.FullProfile{}
	err := p.DB.QueryRow(GetProfileByID, id).Scan(&res.ID, &res.FirstName, &res.LastName, &res.Bio, &res.Avatar)
	if err != nil {
		return nil, fmt.Errorf("get profile by id db: %w", err)
	}
	return res, nil
}

func (p *ProfileRepo) GetAll(self uint32) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.Query(GetAllProfiles, self)
	if err != nil {
		return nil, fmt.Errorf("get all profiles %w", err)
	}
	err = rows.Scan(&res)
	if err != nil {
		return nil, fmt.Errorf("get all profiles %w", err)
	}
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

func (p *ProfileRepo) AddFriendsReq(reciever uint32, sender uint32) error {
	_, err := p.DB.Exec(AddFriends, sender, reciever)
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

func (p *ProfileRepo) GetAllFriends(u uint32) ([]*models.ShortProfile, error) {
	res := make([]*models.ShortProfile, 0)
	rows, err := p.DB.Query(GetAllFriends, u)
	if err != nil {
		return nil, fmt.Errorf("get all friends %w", err)
	}
	err = rows.Scan(&res)
	if err != nil {
		return nil, fmt.Errorf("get all friends %w", err)
	}
	return res, nil
}

func (p *ProfileRepo) CheckStatus(u1 uint32, u2 uint32) (int, error) {
	var status int
	err := p.DB.QueryRow(CheckFriendReq, u1, u2).Scan(&status)
	if err != nil {
		return status, fmt.Errorf("check status db: %w", err)
	}
	return status, nil
}
