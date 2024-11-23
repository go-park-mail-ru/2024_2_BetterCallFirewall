package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type CSATRepository struct {
	DB *sql.DB
}

func NewCSATRepository(db *sql.DB) *CSATRepository {
	return &CSATRepository{
		DB: db,
	}
}

func (cs *CSATRepository) SaveMetrics(ctx context.Context, csat *models.CSAT) error {
	_, err := cs.DB.ExecContext(ctx, InsertNewMetric, csat.InTotal, csat.Review)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CSATRepository) GetMetrics(ctx context.Context, since, before time.Time) (*models.CSATResult, error) {
	res := &models.CSATResult{}
	err := cs.DB.QueryRowContext(ctx, GetMetrics, since, before).Scan(&res.InTotalGrade)
	if err != nil {
		return nil, err
	}
	return res, nil
}
