package errors

var (
	_ error = (*JSONError)(nil)

	// ErrBodyBind returns 429 error.
	ErrBodyBind = JSONError{
		Code:    429,
		Message: "请求异常",
	}

	// ErrNotFound returns 404 error.
	ErrNotFound = JSONError{
		Code:    404,
		Message: "未查询到相关数据",
	}
)

func New(code int, msg string) error {
	return JSONError{
		Code:    code,
		Message: msg,
	}
}

// JSONError details HTTP server error.
type JSONError struct {
	Code    int
	Message string
}

func (j JSONError) Error() string {
	return j.Message
}
