package dto

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CommonErrorResponse struct {
	Status int         `json:"status"`
	Error  *CommonError `json:"error"`
}

type CommonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
