package util

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

func NewErrorResponse(msg, err string) Response {
	return Response{
		Status: _STATUS_FAIL,
		Msg:    msg,
		Error:  err,
	}
}

func NewDataResponse(msg string, data interface{}) Response {
	return Response{
		Status: _STATUS_OK,
		Msg:    msg,
		Data:   data,
	}
}
