package services

import (
	"context"
	"strings"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/constants"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
	sqlite_repository "github.com/Chayakorn2002/pms-classroom-backend/internal/adapters/repositories/sqlite"
	sqlc "github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/sqlc/gen"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/encryption"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/validation"
	"github.com/google/uuid"
	"google.golang.org/api/classroom/v1"
)

type UserService interface {
	RegisterStudent(ctx context.Context, in *dto.RegisterStudentRequest) (*dto.RegisterStudentResponse, error)
	LoginStudent(ctx context.Context, in *dto.LoginStudentRequest) (*dto.LoginStudentResponse, error)
	GetUserProfileByEmail(ctx context.Context, in *dto.GetUserProfileByEmailRequest) (*dto.GetUserProfileByEmailResponse, error)
	// GetUsers(ctx context.Context, in *dto.GetUsersRequest) (*dto.GetUsersResponse, error)
	// GetUserByID(ctx context.Context, in *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error)
	// CreateUser(ctx context.Context, in *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
}

type userService struct {
	userRepo  sqlite_repository.UserRepository
	errors    *exceptions.ApplicationError
	classroom *classroom.Service
}

func NewUserService(
	userRepo sqlite_repository.UserRepository,
	err *exceptions.ApplicationError,
	classroom *classroom.Service,
) UserService {
	return &userService{
		userRepo:  userRepo,
		errors:    err,
		classroom: classroom,
	}
}

func (s *userService) RegisterStudent(ctx context.Context, in *dto.RegisterStudentRequest) (*dto.RegisterStudentResponse, error) {
	if err := validation.ValidateStruct(in); err != nil {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
	}

	email := strings.ToLower(in.Email)
	student, err := s.userRepo.CheckUserExistsByEmail(ctx, email)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	} else if student != nil {
		return nil, s.errors.ErrUserAlreadyExists.WithDebugMessage("email already exists")
	}

	students, err := s.classroom.Courses.Students.List(in.CourseID).Do()
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	var (
		userID    string
		firstname string
		lastname  string
	)
	for _, student := range students.Students {
		if strings.EqualFold(student.Profile.EmailAddress, email) {
			userID = student.UserId
			firstname = student.Profile.Name.GivenName
			lastname = student.Profile.Name.FamilyName
			break
		}
	}
	if userID == "" {
		return nil, s.errors.ErrEmailNotFound.WithDebugMessage("student email not found in the course")
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	password, err := encryption.HashPassword(in.Password)
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	err = s.userRepo.RegisterStudent(ctx, &sqlc.RegisterStudentParams{
		ID:                       uuid.String(),
		Firstname:                firstname,
		Lastname:                 lastname,
		Email:                    email,
		Password:                 password,
		CourseID:                 in.CourseID,
		GoogleClassroomStudentID: userID,
	})
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	return &dto.RegisterStudentResponse{
		Status: constants.StatusSuccess,
	}, nil
}

func (s *userService) LoginStudent(ctx context.Context, in *dto.LoginStudentRequest) (*dto.LoginStudentResponse, error) {
	if err := validation.ValidateStruct(in); err != nil {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
	}

	email := strings.ToLower(in.Email)
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	} else if user == nil {
		return nil, s.errors.ErrInvalidCredential.WithDebugMessage("invalid email")
	}

	err = encryption.CheckPasswordHash(in.Password, user.Password)
	if err != nil {
		return nil, s.errors.ErrInvalidCredential.WithDebugMessage("invalid password")
	}

	return &dto.LoginStudentResponse{
		Status: constants.StatusSuccess,
	}, nil
}

func (s *userService) GetUserProfileByEmail(ctx context.Context, in *dto.GetUserProfileByEmailRequest) (*dto.GetUserProfileByEmailResponse, error) {
	if err := validation.ValidateStruct(in); err != nil {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
	}

	email := strings.ToLower(in.Email)
	user, err := s.userRepo.GetUserProfileByEmail(ctx, email)
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	return &dto.GetUserProfileByEmailResponse{
		Status: constants.StatusSuccess,
		UserProfile: dto.GetUserProfileByEmailResponse_UserProfile{
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			CourseID:  user.CourseID,
			StudentID: user.GoogleClassroomStudentID,
		},
	}, nil
}

// func (s *userService) CreateUser(ctx context.Context, in *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
// 	if err := validation.ValidateStruct(in); err != nil {
// 		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
// 	}

// 	uuid, err := uuid.NewV7()
// 	if err != nil {
// 		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
// 	}

// 	hashedPassword, err := encryption.HashPassword(in.Password)
// 	if err != nil {
// 		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
// 	}

// 	err = s.userRepo.CreateUser(ctx, &sqlc.CreateUserParams{
// 		ID:       uuid.String(),
// 		Email:    in.Email,
// 		Password: hashedPassword,
// 	})
// 	if err != nil {
// 		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
// 	}

// 	return &dto.CreateUserResponse{
// 		Status: constants.StatusSuccess,
// 	}, nil
// }

// func (s *userService) GetUsers(ctx context.Context, in *dto.GetUsersRequest) (*dto.GetUsersResponse, error) {
// 	users, err := s.userRepo.GetUsers(ctx)
// 	if err != nil {
// 		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
// 	}

// 	usersResp := make([]*dto.GetUsersResponseUser, 0)
// 	for _, user := range users {
// 		usersResp = append(usersResp, &dto.GetUsersResponseUser{
// 			ID:        user.ID,
// 			Email:     user.Email,
// 			CreatedAt: user.CreatedAt,
// 			CreatedBy: user.CreatedBy,
// 			UpdatedAt: user.UpdatedAt.Time,
// 			UpdatedBy: user.UpdatedBy.String,
// 		})
// 	}

// 	return &dto.GetUsersResponse{
// 		Status: constants.StatusSuccess,
// 		Data: dto.GetUsersResponseData{
// 			Users: usersResp,
// 		},
// 	}, nil
// }

// func (s *userService) GetUserByID(ctx context.Context, in *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error) {
// 	if err := validation.ValidateStruct(in); err != nil {
// 		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
// 	}

// 	uuid, err := uuid.Parse(in.ID)
// 	if err != nil {
// 		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Error())
// 	}

// 	user, err := s.userRepo.GetUserByID(ctx, uuid.String())
// 	if err != nil {
// 		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
// 	}

// 	return &dto.GetUserByIDResponse{
// 		Status: constants.StatusSuccess,
// 		Data: dto.GetUserByIDResponseData{
// 			ID:        user.ID,
// 			Email:     user.Email,
// 			CreatedAt: user.CreatedAt,
// 			CreatedBy: user.CreatedBy,
// 			UpdatedAt: user.UpdatedAt.Time,
// 			UpdatedBy: user.UpdatedBy.String,
// 		},
// 	}, nil
// }
