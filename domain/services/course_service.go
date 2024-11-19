package services

import (
	"context"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/constants"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/googleapi"
)

type CourseService interface {
	ListCourses(ctx context.Context, in *dto.ListCoursesRequest) (*dto.ListCoursesResponse, error)
}

type courseService struct {
	errors    *exceptions.ApplicationError
	classroom *classroom.Service
}

func NewCourseService(
	err *exceptions.ApplicationError,
	classroom *classroom.Service,
) CourseService {
	return &courseService{
		errors:    err,
		classroom: classroom,
	}
}

func (s *courseService) ListCourses(ctx context.Context, in *dto.ListCoursesRequest) (*dto.ListCoursesResponse, error) {
	var allCourses []*classroom.Course
	pageToken := ""
	for {
		courses, err := s.classroom.Courses.List().Fields(googleapi.Field("courses(id,name,section)")).PageToken(pageToken).Do()
		if err != nil {
			return nil, s.errors.ErrInternal.WithDebugMessage(err.Error())
		}

		allCourses = append(allCourses, courses.Courses...)
		if courses.NextPageToken == "" {
			break
		}
		pageToken = courses.NextPageToken
	}

	resp := make([]dto.ListCoursesResponse_Course, 0, len(allCourses))
	for _, course := range allCourses {
		resp = append(resp, dto.ListCoursesResponse_Course{
			ID:      course.Id,
			Name:    course.Name,
			Section: course.Section,
		})
	}

	return &dto.ListCoursesResponse{
		Status:  constants.StatusSuccess,
		Courses: resp,
	}, nil
}
