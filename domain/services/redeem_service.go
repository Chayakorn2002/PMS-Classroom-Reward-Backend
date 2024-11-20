package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
	"github.com/Chayakorn2002/pms-classroom-backend/internal/adapters/repositories/sqlite"
	sqlc "github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/sqlc/gen"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/gen"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/validation"
	"github.com/google/uuid"
	"google.golang.org/api/classroom/v1"
)

type RedeemService interface {
	RedeemReward(ctx context.Context, in *dto.RedeemRewardRequest) (*dto.RedeemRewardResponse, error)
}

type redeemService struct {
	redeemLogRepo sqlite.RedeemLogRepository
	errors        *exceptions.ApplicationError
	classroom     *classroom.Service
}

func NewRedeemService(
	redeemLogRepo sqlite.RedeemLogRepository,
	err *exceptions.ApplicationError,
	classroom *classroom.Service,
) RedeemService {
	return &redeemService{
		errors:        err,
		classroom:     classroom,
		redeemLogRepo: redeemLogRepo,
	}
}

func (s *redeemService) RedeemReward(ctx context.Context, in *dto.RedeemRewardRequest) (*dto.RedeemRewardResponse, error) {
	fmt.Println("in", in)
	if err := validation.ValidateStruct(in); err != nil {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
	}

	assignment, err := s.classroom.Courses.CourseWork.Get(in.CourseID, in.AssignmentID).Do()
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	// Check if the assignment is already redeemed
	redeemLog, err := s.redeemLogRepo.GetRedeemLogsByStudentId(ctx, in.StudentID)
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}
	for _, log := range redeemLog {
		if log.AssignmentID == assignment.Id {
			return nil, s.errors.ErrAlreadyRedeeemed.WithDebugMessage("The assignment is already redeemed")
		}
	}

	submissions, err := s.classroom.Courses.CourseWork.StudentSubmissions.List(in.CourseID, in.AssignmentID).UserId(in.StudentID).Do()
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	if len(submissions.StudentSubmissions) == 0 {
		return nil, s.errors.ErrBadRequest.WithDebugMessage("No submission found")
	}

	// Get the latest submission
	latestSubmission := submissions.StudentSubmissions[0]
	for _, submission := range submissions.StudentSubmissions {
		submissionTime, err := time.Parse(time.RFC3339, submission.CreationTime)
		if err != nil {
			return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
		}
		latestSubmissionTime, err := time.Parse(time.RFC3339, latestSubmission.CreationTime)
		if err != nil {
			return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
		}
		if submissionTime.After(latestSubmissionTime) {
			latestSubmission = submission
		}
	}

	if latestSubmission.State != "RETURNED" {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(fmt.Sprintf("The submission state is %s", latestSubmission.State))
	}

	// Check if the student has passed the assignment
	// 80% of the max points
	if latestSubmission.AssignedGrade < assignment.MaxPoints*0.8 {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(fmt.Sprintf("The student has not passed the assignment. Assigned grade: %f", latestSubmission.AssignedGrade))
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	// Generate Serial Number for the reward in 10 alphanumeric characters
	serial := gen.GenerateSerial(10)

	// Redeem the assignment
	err = s.redeemLogRepo.CreateRedeemLog(ctx, &sqlc.CreateRedeemLogParams{
		ID:                       uuid.String(),
		Serial:                   serial,
		CourseID:                 in.CourseID,
		GoogleClassroomStudentID: in.StudentID,
		AssignmentID:             in.AssignmentID,
	})
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	return &dto.RedeemRewardResponse{
		Status: 1000,
		Serial: serial,
	}, nil
}
