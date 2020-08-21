package out

var (
	_ error = (*Bucket)(nil)

	// ErrBodyBind returns 429 error.
	ErrBodyBind = Bucket{
		Code:    429,
		Message: "请求异常",
	}

	// ErrNotFound returns 404 error.
	ErrNotFound = Bucket{
		Code:    404,
		Message: "未查询到相关数据",
	}

	// Err500 returns 500 error.
	Err500 = Bucket{
		Code:    500,
		Message: "💥服务器爆炸啦",
	}

	// Err401 returns 401 error.
	Err401 = Bucket{
		Code:    401,
		Message: "请登录",
	}
)

// New error with given http status code and message.
func New(code int, msg string, data interface{}) Bucket {
	return Bucket{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

// Data wraps data into bucket.
func Data(data interface{}) Bucket {
	return Bucket{
		Code:    0,
		Message: "OK",
		Data:    data,
	}
}

// Err wraps error message into bucket.
func Err(code int, msg string) Bucket {
	return Bucket{
		Code:    code,
		Message: msg,
	}
}

// OR returns leftValue when cond is true, or returns rightValue.
func OR(cond bool, leftValue interface{}, rightValue interface{}) interface{} {
	if cond {
		return leftValue
	}
	return rightValue
}

// Bucket details HTTP server error.
type Bucket struct {
	Code    int
	Message string      `json:",omitempty"`
	Data    interface{} `json:",omitempty"`
}

func (b Bucket) Error() string {
	return b.Message
}
