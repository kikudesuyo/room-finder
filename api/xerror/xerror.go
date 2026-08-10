package xerror

import (
	"net/http"
	"strings"

	"github.com/morikuni/failure"
)

func GetHTTPErrorCode(err error) int {
	code, ok := failure.CodeOf(err)
	if !ok {
		return http.StatusInternalServerError
	}
	switch {
	case strings.HasPrefix(code.ErrorCode(), "A-"):
		return http.StatusBadRequest
	case strings.HasPrefix(code.ErrorCode(), "E-"):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func GetStringCode(err error) string {
	code, ok := failure.CodeOf(err)
	if !ok {
		return "E-0-00001"
	}
	return code.ErrorCode()
}

func GetErrorMessage(err error) string {
	message, ok := failure.MessageOf(err)
	if !ok {
		return err.Error()
	}
	return message
}
