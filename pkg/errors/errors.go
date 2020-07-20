package errors

import "fmt"

var (
	_ error = (*JSONError)(nil)

	// ErrBodyBind returns 429 error.
	ErrBodyBind = JSONError{
		Code: 429,
		Msg:  "请求异常",
	}

	// ErrNotFound returns 404 error.
	ErrNotFound = JSONError{
		Code: 404,
		Msg:  "未查询到相关数据",
	}
)

// JSONError details HTTP server error.
type JSONError struct {
	Code int
	Msg  string
}

func (j JSONError) Error() string {
	return fmt.Sprintf("HTTP code %d caused by %s", j.Code, j.Msg)
}
