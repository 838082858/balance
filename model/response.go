package model

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewErrResponse(code int, errMsg string, model interface{}) Response {
	return Response{
		Code:    code,
		Message: errMsg,
		Data:    model,
	}
}
