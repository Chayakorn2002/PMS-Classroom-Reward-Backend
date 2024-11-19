package dto

type ListStudentAssignmentRequest struct {
	StudentID string `json:"student_id" validate:"required"`
	CourseID  string `json:"course_id" validate:"required"`
}

type ListStudentAssignmentResponse struct {
	Status      int                                        `json:"status"`
	Assignments []ListStudentAssignmentResponse_Assignment `json:"assignments"`
}

type ListStudentAssignmentResponse_Assignment struct {
	AssignmentID  string  `json:"assignment_id"`
	Title         string  `json:"title"`
	AssignedGrade float64 `json:"assigned_grade"`
	MaxPoints     float64 `json:"max_points"`
	State         string  `json:"state"`
	CreationTime  string  `json:"creation_time"`
}
