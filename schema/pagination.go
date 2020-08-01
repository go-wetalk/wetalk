package schema

// Pagination 分页响应结构体
type Pagination struct {
	RowCount int
	PerPage  int
	Data     interface{}
}
