package services

import (
	"github.com/Chayakorn2002/pms-classroom-backend/config"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
	sqlite_repository "github.com/Chayakorn2002/pms-classroom-backend/internal/adapters/repositories/sqlite"
	"google.golang.org/api/classroom/v1"
)

type Service struct {
	config    *config.Config
	errors    *exceptions.ApplicationError
	classroom *classroom.Service

	UserService       UserService
	CourseService     CourseService
	AssignmentService AssignmentService
	RedeemService     RedeemService
}

func NewService(
	repo *sqlite_repository.SqliteRepository,
	cfg *config.Config,
	errors *exceptions.ApplicationError,
	classroom *classroom.Service,
) Service {
	return Service{
		config:            cfg,
		errors:            errors,
		classroom:         classroom,
		UserService:       NewUserService(repo.UserRepository, errors, classroom),
		CourseService:     NewCourseService(errors, classroom),
		AssignmentService: NewAssignmentService(repo.RedeemLogRepository, errors, classroom),
		RedeemService:     NewRedeemService(repo.RedeemLogRepository, errors, classroom),
	}
}
