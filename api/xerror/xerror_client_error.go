package xerror

import (
	"errors"

	"github.com/morikuni/failure"
)

const (
	CodeClientValidation       failure.StringCode = "A-1-10001"
	CodeClientResourceNotFound failure.StringCode = "A-0-00001"
)

func ClientValidationErr(err error, meta ...map[string]string) error {
	if err == nil {
		err = errors.New("request validation failed")
	}
	return failure.Translate(err, CodeClientValidation, failure.Message("request validation failed"), failure.Context(metaValue(meta)))
}

func ClientResourceNotFoundErr(meta ...map[string]string) error {
	return failure.New(CodeClientResourceNotFound, failure.Message("resource not found"), failure.Context(metaValue(meta)))
}

func metaValue(meta []map[string]string) map[string]string {
	if len(meta) == 0 || meta[0] == nil {
		return map[string]string{}
	}
	return meta[0]
}
