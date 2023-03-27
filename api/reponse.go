package api

type ReturnResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *DataResult `json:"data,omitempty"`
}

type DataResult struct {
	RegisterCode string `json:"register_code,omitempty"`
	Status       string `json:"status,omitempty"`
}

func Err(code int, message string) ReturnResponse {
	return ReturnResponse{
		Code:    code,
		Message: message,
	}
}
