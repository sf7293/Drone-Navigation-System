package domain

type RouterRequestCalculateLocation struct {
	XCord    string `json:"x" form:"x"`
	YCord    string `json:"y" form:"y"`
	ZCord    string `json:"z" form:"z"`
	Velocity string `json:"vel" form:"vel"`
}

type RouterResponseCalculateLocation struct {
	Location float64 `json:"loc"`
}

type ErrorResponseData struct {
	Label           string                 `json:"label"`
	Title           string                 `json:"title"`
	Message         string                 `json:"message"`
	ServerError     string                 `json:"server_error,omitempty"`
	ServerErrorData map[string]interface{} `json:"server_error_data,omitempty"`
	PublicErrorData map[string]interface{} `json:"error_data,omitempty"`
}

type RouterResponseError struct {
	Success bool              `json:"success"`
	Data    ErrorResponseData `json:"data"`
}
