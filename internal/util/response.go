package util

import "fmt"

const (
	_STATUS_OK   = "OK"
	_STATUS_FAIL = "FAIL"
)

type Response struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func NewMsgResponse(msg string) Response {
	return Response{
		Status: _STATUS_OK,
		Msg:    msg,
	}
}

func NewErrorResponse(msg string, err error) Response {
	if err != nil {
		return Response{
			Status: _STATUS_FAIL,
			Msg:    msg,
			Error:  fmt.Sprintf("%s", err),
		}
	}
	return Response{
		Status: _STATUS_FAIL,
		Msg:    msg,
	}
}

func NewDataResponse(msg string, data interface{}) Response {
	return Response{
		Status: _STATUS_OK,
		Msg:    msg,
		Data:   data,
	}
}
