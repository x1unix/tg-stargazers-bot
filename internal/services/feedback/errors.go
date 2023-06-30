package feedback

import "fmt"

type ErrorCode uint

const (
	ErrNoError             ErrorCode = 0
	ErrTokenSaveError      ErrorCode = 0xf1
	ErrBadAuthCallbackCall ErrorCode = 0xf2
	ErrInvalidToken        ErrorCode = 0xf100
	ErrUnknown             ErrorCode = 0xffffff
)

func NewErrorMessage(msg string, code ErrorCode) string {
	return fmt.Sprintf("😬 %s, please try again later.\nError code: `%06X`", msg, code)
}
