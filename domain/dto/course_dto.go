package dto

type ListCoursesRequest struct{}

type ListCoursesResponse struct {
	Status  int                          `json:"status"`
	Courses []ListCoursesResponse_Course `json:"courses"`
}

type ListCoursesResponse_Course struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Section string `json:"section"`
}
