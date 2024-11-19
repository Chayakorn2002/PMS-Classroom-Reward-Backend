package server

import (
	"context"
	"net/http"

	"github.com/Chayakorn2002/pms-classroom-backend/domain/dto"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/services"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/not_found"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/router"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/transport"
)

func registerRoute(service services.Service) http.Handler {
	mux := http.NewServeMux()
	r := router.NewRouter(mux)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		not_found.NotFound(w, r)
	}))

	r.Post("/api/v1/health-check",
		transport.NewTransport(
			&dto.HealthCheckRequest{},
			transport.NewEndpoint(func(ctx context.Context, in *dto.HealthCheckRequest) (*dto.HealthCheckResponse, error) {
				return &dto.HealthCheckResponse{
					Status:  1000,
					Message: "OK",
				}, nil
			})))

	// api
	{
		// v1
		{
			// courses
			{
				r.Post("/api/v1/courses/list-courses",
					transport.NewTransport(
						&dto.ListCoursesRequest{},
						transport.NewEndpoint(
							service.CourseService.ListCourses,
						)))
			}

			// assignments
			{
				r.Post("/api/v1/assignments/list-student-assignment",
					transport.NewTransport(
						&dto.ListStudentAssignmentRequest{},
						transport.NewEndpoint(
							service.AssignmentService.ListStudentAssignment,
						)))
			}

			// redeem
			{
				r.Post("/api/v1/redeem/redeem-reward",
					transport.NewTransport(
						&dto.RedeemRewardRequest{},
						transport.NewEndpoint(
							service.RedeemService.RedeemReward,
						)))
			}

			// users
			{
				r.Post("/api/v1/users/login-student",
					transport.NewTransport(
						&dto.LoginStudentRequest{},
						transport.NewEndpoint(
							service.UserService.LoginStudent,
						)))

				r.Post("/api/v1/users/register-student",
					transport.NewTransport(
						&dto.RegisterStudentRequest{},
						transport.NewEndpoint(
							service.UserService.RegisterStudent,
						)))

				r.Post("/api/v1/users/get-user-profile-by-email",
					transport.NewTransport(
						&dto.GetUserProfileByEmailRequest{},
						transport.NewEndpoint(
							service.UserService.GetUserProfileByEmail,
						)))

				// r.Post("/api/v1/users/get-users",
				// 	transport.NewTransport(
				// 		&dto.GetUsersRequest{},
				// 		transport.NewEndpoint(
				// 			service.UserService.GetUsers,
				// 		)))

				// r.Post("/api/v1/users/get-user-by-id",
				// 	transport.NewTransport(
				// 		&dto.GetUserByIDRequest{},
				// 		transport.NewEndpoint(
				// 			service.UserService.GetUserByID,
				// 		)))

				// r.Post("/api/v1/users/create-user",
				// 	transport.NewTransport(
				// 		&dto.CreateUserRequest{},
				// 		transport.NewEndpoint(
				// 			service.UserService.CreateUser,
				// 		)))
			}
		}
	}

	return r
}
