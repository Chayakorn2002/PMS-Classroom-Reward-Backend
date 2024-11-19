package dto

type RegisterStudentRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	CourseID string `json:"course_id" validate:"required"`
}

type RegisterStudentResponse struct {
	Status int `json:"status"`
}

type LoginStudentRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginStudentResponse struct {
	Status int `json:"status"`
}

type GetUserProfileByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type GetUserProfileByEmailResponse struct {
	Status      int         `json:"status"`
	UserProfile GetUserProfileByEmailResponse_UserProfile `json:"user_profile"`
}

type GetUserProfileByEmailResponse_UserProfile struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"google_classroom_student_id"`
}

// // GetUsers
// type GetUsersRequest struct {
// }

// type GetUsersResponse struct {
// 	Status int                  `json:"status"`
// 	Data   GetUsersResponseData `json:"data"`
// }

// type GetUsersResponseData struct {
// 	Users []*GetUsersResponseUser `json:"users"`
// }

// type GetUsersResponseUser struct {
// 	ID        string    `json:"id"`
// 	Email     string    `json:"email"`
// 	CreatedAt time.Time `json:"created_at"`
// 	CreatedBy string    `json:"created_by"`
// 	UpdatedAt time.Time `json:"updated_at"`
// 	UpdatedBy string    `json:"updated_by"`
// }

// // GetUserByID
// type GetUserByIDRequest struct {
// 	ID string `json:"id" validate:"required"`
// }

// type GetUserByIDResponse struct {
// 	Status int                     `json:"status"`
// 	Data   GetUserByIDResponseData `json:"data"`
// }

// type GetUserByIDResponseData struct {
// 	ID        string    `json:"id"`
// 	Email     string    `json:"email"`
// 	CreatedAt time.Time `json:"created_at"`
// 	CreatedBy string    `json:"created_by"`
// 	UpdatedAt time.Time `json:"updated_at"`
// 	UpdatedBy string    `json:"updated_by"`
// }

// // CreateUser
// type CreateUserRequest struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required"`
// }

// type CreateUserResponse struct {
// 	Status int `json:"status"`
// }
