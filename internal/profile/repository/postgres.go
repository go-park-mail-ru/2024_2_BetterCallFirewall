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

func NewProfileRepo(db *pgx.ConnPool) profile.ProfileUsecase {
	repo := &ProfileRepo{
		DB: db,
	}
	return repo
}

const (
	GetProfileByID  = "SELECT pr.id, per.first_name, per.last_name, pr.bio, pr.avatar FROM profile AS pr INNER JOIN person AS pe ON pe.id = pr.person_id WHERE pr.id = $1;"
	GetProfilePics  = "SELECT pi.image FROM profile AS p INNER JOIN profile_image AS pi ON p.id = pi.profile_id WHERE p.id = $1;"
	GetProfilePosts = "SELECT c.text FROM post AS p INNER JOIN content AS c ON p.content_id = c.id WHERE p.author_id = $1;"
	GetAllProfiles  = "SELECT id, first_name, last_name, avatar FROM person WHERE id != $1 ;"
	UpdateProfile   = "UPDATE profile SET first_name = $1, last_name = $2, bio = $3, avatar = $4 WHERE id = $5;"
	DeleteProfile   = "DELETE FROM profile WHERE id = $1;"
	CreateProfile   = "INSERT INTO profile (first_name, last_name, bio, avatar) VALUES ($1, $2, $3, $4) RETURNING id;);"
)

func (p *ProfileRepo) GetProfileById(id uint32) (*models.Profile, error) {
	res := &models.Profile{}
	err := p.DB.QueryRow(GetProfileByID, id).Scan(&res.ID, &res.FirstName, &res.LastName, &res.Bio, &res.Avatar)
	if err != nil {
		return nil, fmt.Errorf("get profile by id %w", err)
	}
	rows, err := p.DB.Query(GetProfilePosts, id)
	if err != nil {
		return nil, fmt.Errorf("get profile posts %w", err)
	}
	err = rows.Scan(&res.Posts)
	if err != nil {
		return nil, fmt.Errorf("get profile posts %w", err)
	}
	return res, nil
}

func (p *ProfileRepo) GetAll(self uint32) ([]*models.Profile, error) {
	res := make([]*models.Profile, 0)
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

/*func (p *ProfileRepo) CreateProfile(profile models.Profile) (uint32, error) {
	p.DB.Exec() //Возможно через триггеры при создании person
	panic("implement me")
}*/

func (p *ProfileRepo) UpdateProfile(profile *models.Profile) error {
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
