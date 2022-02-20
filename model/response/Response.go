package model

type (
	PageResponse struct {
		Status  int
		Message string
		Items   interface{}
	}

	ErrorResponse struct {
		Status  int
		Message string
	}

	DBHealthResponse struct {
		DBStatus int
		Message  string
	}
)
