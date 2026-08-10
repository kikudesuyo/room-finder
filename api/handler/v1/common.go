package v1

import "strconv"

func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
