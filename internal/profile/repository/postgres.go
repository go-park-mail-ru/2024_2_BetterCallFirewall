package repository

import (
	"github.com/2024_2_BetterCallFirewall/internal/profile/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"github.com/jackc/pgx"
)

type ProfileRepo struct {
	DB pgx.ConnPool
}

func NewProfileRepo() service.ProfileUsecase {}

const (
	GetProfileByID = "SELECT pr.id, per.first_name, per.last_name, pr.bio, pr.avatar FROM profile AS pr INNER JOIN person AS pe ON pe.id = pr.person_id WHERE pr.id = $1;"
	`GetProfilePics = "SELECT pi.image FROM profile AS p INNER JOIN profile_image AS pi ON p.id = pi.profile_id WHERE p.id = $1;"`
	GetProfilePosts = "SELECT c.text FROM post AS p INNER JOIN content AS c INNER JOIN ON p.content_id = c.id WHERE p.author_id = $1;"
	GetAllProfiles  = "SELECT pr.id, per.first_name, per.last_name, pr.bio, pr.avatar FROM profile AS pr INNER JOIN person AS pe ON pe.id = pr.person_id;"
	UpdateProfile   = "UPDATE profile SET first_name = $1, last_name = $2, bio = $3,avatar = $4 WHERE id = $4;"
	DeleteProfile   = "DELETE FROM profile WHERE id = $1;"
	CreateProfile   = "INSERT INTO profile (first_name, last_name, bio, avatar) VALUES ($1, $2, $3, $4) RETURNING id;);"
)

func (p ProfileRepo) GetProfileById(uint64 uint32) (*models.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepo) GetAll() ([]*models.Profile, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepo) CreateProfile(profile models.Profile) (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepo) UpdateProfile(profile *models.Profile) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProfileRepo) DeleteProfile(u uint32) error {
	//TODO implement me
	panic("implement me")
}
