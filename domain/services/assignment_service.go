package services

import (
	"context"
	"sync"
	"time"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
	"github.com/Chayakorn2002/pms-classroom-backend/internal/adapters/repositories/sqlite"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/validation"
	"google.golang.org/api/classroom/v1"
)

type AssignmentService interface {
	ListStudentAssignment(ctx context.Context, in *dto.ListStudentAssignmentRequest) (*dto.ListStudentAssignmentResponse, error)
}

type assignmentService struct {
	redeemLogRepo sqlite.RedeemLogRepository
	errors        *exceptions.ApplicationError
	classroom     *classroom.Service
}

func NewAssignmentService(
	redeemLogRepo sqlite.RedeemLogRepository,
	err *exceptions.ApplicationError,
	classroom *classroom.Service,
) AssignmentService {
	return &assignmentService{
		errors:        err,
		classroom:     classroom,
		redeemLogRepo: redeemLogRepo,
	}
}

func (s *assignmentService) ListStudentAssignment(ctx context.Context, in *dto.ListStudentAssignmentRequest) (*dto.ListStudentAssignmentResponse, error) {
	if err := validation.ValidateStruct(in); err != nil {
		return nil, s.errors.ErrBadRequest.WithDebugMessage(err.Message)
	}

	assignments, err := s.classroom.Courses.CourseWork.List(in.CourseID).Fields("courseWork(id,title,maxPoints)").Do()
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	redeemLogs, err := s.redeemLogRepo.GetRedeemLogsByStudentId(ctx, in.StudentID)
	if err != nil {
		return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
	}

	isRedeemed := func(assignmentID string) bool {
		for _, log := range redeemLogs {
			if log.AssignmentID == assignmentID {
				return true
			}
		}
		return false
	}

	resp := make([]dto.ListStudentAssignmentResponse_Assignment, 0, len(assignments.CourseWork))
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(assignments.CourseWork))

	for _, courseWork := range assignments.CourseWork {
		wg.Add(1)
		go func(courseWork *classroom.CourseWork) {
			defer wg.Done()
			submissions, err := s.classroom.Courses.CourseWork.StudentSubmissions.List(in.CourseID, courseWork.Id).UserId(in.StudentID).States("RETURNED").Fields("studentSubmissions(id,assignedGrade,state,creationTime)").Do()
			if err != nil {
				errChan <- s.errors.ErrInternal.WithDebugMessage(err.Error())
				return
			}

			if len(submissions.StudentSubmissions) == 0 {
				return
			}

			// Get the latest submission
			latestSubmission := submissions.StudentSubmissions[0]
			for _, submission := range submissions.StudentSubmissions {
				submissionTime, err := time.Parse(time.RFC3339, submission.CreationTime)
				if err != nil {
					errChan <- s.errors.ErrInternal.WithDebugMessage(err.Error())
					return
				}
				latestSubmissionTime, err := time.Parse(time.RFC3339, latestSubmission.CreationTime)
				if err != nil {
					errChan <- s.errors.ErrInternal.WithDebugMessage(err.Error())
					return
				}
				if submissionTime.After(latestSubmissionTime) {
					latestSubmission = submission
				}
			}

			mu.Lock()
			defer mu.Unlock()

			// Check if the assignment is redeemed
			if isRedeemed(courseWork.Id) {
				resp = append(resp, dto.ListStudentAssignmentResponse_Assignment{
					AssignmentID:  courseWork.Id,
					Title:         courseWork.Title,
					AssignedGrade: latestSubmission.AssignedGrade,
					MaxPoints:     courseWork.MaxPoints,
					State:         "REDEEMED",
					CreationTime:  latestSubmission.CreationTime,
				})
				return
			}

			resp = append(resp, dto.ListStudentAssignmentResponse_Assignment{
				AssignmentID:  courseWork.Id,
				Title:         courseWork.Title,
				AssignedGrade: latestSubmission.AssignedGrade,
				MaxPoints:     courseWork.MaxPoints,
				State:         latestSubmission.State,
				CreationTime:  latestSubmission.CreationTime,
			})
		}(courseWork)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return &dto.ListStudentAssignmentResponse{
		Status:      1000,
		Assignments: resp,
	}, nil
}
