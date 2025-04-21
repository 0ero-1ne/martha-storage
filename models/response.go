package models

type Response struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}

func NewErrorResponse(err string) Response {
	return Response{
		Data:  nil,
		Error: err,
	}
}

func NewSuccessResponse(data any) Response {
	return Response{
		Data:  data,
		Error: nil,
	}
}
