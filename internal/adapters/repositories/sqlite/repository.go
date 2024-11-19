package sqlite

import (
	"context"
	"fmt"

	"github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/sqlite"
)

type SqliteRepository struct {
	UserRepository      UserRepository
	RedeemLogRepository RedeemLogRepository
}

func NewSqliteRepository(ctx context.Context) (*SqliteRepository, error) {
	db, err := sqlite.OpenSQLiteDB(ctx)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	return &SqliteRepository{
		UserRepository:      NewUserRepository(db),
		RedeemLogRepository: NewRedeemLogRepository(db),
	}, nil
}
