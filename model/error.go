package model

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
