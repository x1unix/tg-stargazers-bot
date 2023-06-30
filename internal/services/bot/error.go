package bot

import (
	"errors"
	"fmt"
)

// ErrorResponse is error response that will be sent to a client.
type ErrorResponse struct {
	msg string
	err error
}

func NewErrorResponse(msg string, err error) *ErrorResponse {
	return &ErrorResponse{msg: msg, err: err}
}

func (err ErrorResponse) Error() string {
	if err.msg == "" {
		return err.err.Error()
	}

	if err.err == nil {
		return err.msg
	}

	return fmt.Sprint(err.msg, ": ", err.err.Error())
}

func IsErrorResponse(err error) (*ErrorResponse, bool) {
	rsp := new(ErrorResponse)
	ok := errors.As(err, &rsp)
	if !ok {
		return nil, false
	}

	return rsp, true
}
