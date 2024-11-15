package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const LIMIT = 20

type ProfileRepo struct {
	DB *sql.DB
}

func NewProfileRepo(db *sql.DB) *ProfileRepo {
	repo := &ProfileRepo{
		DB: db,
	}
	return repo
}

func (p *ProfileRepo) Create(user *models.User, ctx context.Context) (uint32, error) {
	var id uint32
	err := p.DB.QueryRowContext(ctx, CreateUser, user.FirstName, user.LastName, user.Email, user.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("postgres create user: %w", my_err.ErrUserAlreadyExists)
		}
		return 0, fmt.Errorf("postgres create user: %w", err)
	}

	return id, nil
}

func (p *ProfileRepo) GetByEmail(email string, ctx context.Context) (*models.User, error) {
	user := &models.User{}
	err := p.DB.QueryRowContext(ctx, GetUserByEmail, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("postgres get user: %w", my_err.ErrUserNotFound)
		}
		return nil, fmt.Errorf("postgres get user: %w", err)
	}

	return user, nil
}

func (p *ProfileRepo) GetProfileById(ctx context.Context, id uint32) (*models.FullProfile, error) {
	res := &models.FullProfile{}
	err := p.DB.QueryRowContext(ctx, GetProfileByID, id).Scan(&res.ID, &res.FirstName, &res.LastName, &res.Bio, &res.Avatar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrProfileNotFound
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

func (p *ProfileRepo) GetStatuses(ctx context.Context, self uint32) ([]uint32, []uint32, []uint32, error) {
	var (
		friends          []uint32
		subscribers      []uint32
		subscriptions    []uint32
		tmpFriends       sql.NullString
		tmpSubscribers   sql.NullString
		tmpSubscriptions sql.NullString
	)

	err := p.DB.QueryRowContext(ctx, GetAllStatuses, self).Scan(&tmpFriends, &tmpSubscribers, &tmpSubscriptions)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get all statuses query: %w", err)
	}
	if tmpFriends.Valid {
		err = json.Unmarshal([]byte(tmpFriends.String), &friends)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get all statuses json parsing: %w", err)
		}
	}
	if tmpSubscribers.Valid {
		err = json.Unmarshal([]byte(tmpSubscribers.String), &subscribers)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get all statuses json parsing: %w", err)
		}
	}
	if tmpSubscriptions.Valid {
		err = json.Unmarshal([]byte(tmpSubscriptions.String), &subscriptions)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("get all statuses json parsing: %w", err)
		}
	}

	return friends, subscribers, subscriptions, nil

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

func (p *ProfileRepo) UpdateProfile(ctx context.Context, profile *models.FullProfile) error {
	_, err := p.DB.ExecContext(ctx, UpdateProfile, profile.FirstName, profile.LastName, profile.Bio, profile.ID)
	if err != nil {
		return fmt.Errorf("update profile %w", err)
	}

	return nil
}

func (p *ProfileRepo) UpdateWithAvatar(ctx context.Context, newProfile *models.FullProfile) error {
	_, err := p.DB.ExecContext(ctx, UpdateProfileAvatar, newProfile.ID, newProfile.Avatar, newProfile.FirstName, newProfile.LastName, newProfile.Bio)
	if err != nil {
		return fmt.Errorf("update profile with avatar %w", err)
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

func (p *ProfileRepo) CheckFriendship(ctx context.Context, self uint32, profile uint32) (bool, error) {
	_, err := p.DB.ExecContext(ctx, CheckFriendship, self, profile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("check friendship: %w", err)
	}
	return false, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, my_err.ErrProfileNotFound
		}
		return nil, fmt.Errorf("get header db: %w", err)
	}
	return profile, nil
}

func (p *ProfileRepo) GetCommunitySubs(ctx context.Context, communityID uint32, lastInsertId uint32) ([]*models.ShortProfile, error) {
	var subs []*models.ShortProfile
	rows, err := p.DB.QueryContext(ctx, GetCommunitySubs, communityID, lastInsertId, LIMIT)
	if err != nil {
		return nil, fmt.Errorf("get community subs db: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		profile := &models.ShortProfile{}
		err = rows.Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Avatar)
		if err != nil {
			return nil, fmt.Errorf("get community subs db: %w", err)
		}
		subs = append(subs, profile)
	}
	return subs, nil
}
