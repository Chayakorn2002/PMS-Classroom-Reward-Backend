package sqlite

import (
	"context"
	"database/sql"

	sqlc "github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/sqlc/gen"
)

type UserRepository interface {
	CheckUserExistsByEmail(ctx context.Context, email string) (*sqlc.CheckUserExistsByEmailRow, error)
	RegisterStudent(ctx context.Context, in *sqlc.RegisterStudentParams) error
	GetUserByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error)
	GetUserProfileByEmail(ctx context.Context, email string) (*sqlc.GetUserProfileByEmailRow, error)
	// GetUsers(ctx context.Context) ([]sqlc.GetUsersRow, error)
	// GetUserByID(ctx context.Context, id string) (sqlc.GetUserByIdRow, error)
	// CreateUser(ctx context.Context, in *sqlc.CreateUserParams) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CheckUserExistsByEmail(ctx context.Context, email string) (*sqlc.CheckUserExistsByEmailRow, error) {
	qtx := sqlc.New(r.db)
	user, err := qtx.CheckUserExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) RegisterStudent(ctx context.Context, in *sqlc.RegisterStudentParams) error {
	qtx := sqlc.New(r.db)
	err := qtx.RegisterStudent(ctx, *in)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error) {
	qtx := sqlc.New(r.db)
	user, err := qtx.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserProfileByEmail(ctx context.Context, email string) (*sqlc.GetUserProfileByEmailRow, error) {
	qtx := sqlc.New(r.db)
	user, err := qtx.GetUserProfileByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// func (r *userRepository) GetUsers(ctx context.Context) ([]sqlc.GetUsersRow, error) {
// 	qtx := sqlc.New(r.db)
// 	users, err := qtx.GetUsers(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }

// func (r *userRepository) GetUserByID(ctx context.Context, id string) (sqlc.GetUserByIdRow, error) {
// 	qtx := sqlc.New(r.db)
// 	user, err := qtx.GetUserById(ctx, id)
// 	if err != nil {
// 		return sqlc.GetUserByIdRow{}, err
// 	}

// 	return user, nil
// }

// func (r *userRepository) CreateUser(ctx context.Context, in *sqlc.CreateUserParams) error {
// 	qtx := sqlc.New(r.db)
// 	err := qtx.CreateUser(ctx, *in)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
