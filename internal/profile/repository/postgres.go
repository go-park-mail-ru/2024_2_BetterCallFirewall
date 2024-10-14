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
		return nil, fmt.Errorf("get profile by id %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("get profile posts %w", err)
	}
	err = rows.Scan(&res.Posts)
	if err != nil {
		return nil, fmt.Errorf("get profile posts %w", err)
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

/*func (p *ProfileRepo) CreateProfile(profile models.FullProfile) (uint32, error) {
	p.DB.Exec() //Возможно через триггеры при создании person
	panic("implement me")
}*/

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
