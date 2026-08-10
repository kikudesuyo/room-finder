package xerror

import "github.com/morikuni/failure"

const CodeDatabase failure.StringCode = "E-0-10001"

func UnknownDBErr(err error, meta ...map[string]string) error {
	if err == nil {
		return nil
	}
	return failure.Translate(err, CodeDatabase, failure.Message("database error"), failure.Context(metaValue(meta)))
}
