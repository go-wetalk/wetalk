package errors

type JSONError struct {
	Code int
	Msg  string
}

var (
	ErrBodyBind = JSONError{
		Code: 429,
		Msg:  "请求异常",
	}

	ErrNotFound = JSONError{
		Code: 404,
		Msg:  "未查询到相关数据",
	}
)
