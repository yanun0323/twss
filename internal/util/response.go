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

func NewErrorResponse(msg string, errs ...error) Response {
	if len(errs) == 0 {
		return Response{
			Status: _STATUS_FAIL,
			Msg:    msg,
		}
	}
	return Response{
		Status: _STATUS_FAIL,
		Msg:    msg,
		Error:  fmt.Sprintf("%s", errs[0]),
	}
}

func NewDataResponse(msg string, data interface{}) Response {
	return Response{
		Status: _STATUS_OK,
		Msg:    msg,
		Data:   data,
	}
}
