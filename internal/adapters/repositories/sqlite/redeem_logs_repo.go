package sqlite

import (
	"context"
	"database/sql"

	sqlc "github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/sqlc/gen"
)

type RedeemLogRepository interface {
	GetRedeemLogsByStudentId(ctx context.Context, studentID string) ([]sqlc.GetRedeemLogsByStudentIdRow, error)
	CreateRedeemLog(ctx context.Context, in *sqlc.CreateRedeemLogParams) error
}

type redeemLogRepository struct {
	db *sql.DB
}

func NewRedeemLogRepository(db *sql.DB) RedeemLogRepository {
	return &redeemLogRepository{
		db: db,
	}
}

func (r *redeemLogRepository) GetRedeemLogsByStudentId(ctx context.Context, studentID string) ([]sqlc.GetRedeemLogsByStudentIdRow, error) {
	qtx := sqlc.New(r.db)
	rows, err := qtx.GetRedeemLogsByStudentId(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *redeemLogRepository) CreateRedeemLog(ctx context.Context, in *sqlc.CreateRedeemLogParams) error {
	qtx := sqlc.New(r.db)
	err := qtx.CreateRedeemLog(ctx, *in)
	if err != nil {
		return err
	}

	return nil
}
