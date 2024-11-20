package dto

type RedeemRewardRequest struct {
	StudentID    string `json:"student_id" validate:"required"`
	CourseID     string `json:"course_id" validate:"required"`
	AssignmentID string `json:"assignment_id" validate:"required"`
}

type RedeemRewardResponse struct {
	Status int    `json:"status"`
	Serial string `json:"serial"`
}
