package schema

// Pagination 分页响应结构体
type Pagination struct {
	RowCount int
	PerPage  int
	Data     interface{}
}

type Paginate struct {
	Page int
	Size int
}

func (p Paginate) Offset() int {
	if p.Page < 1 {
		return 0
	}

	return (p.Page - 1) * p.Size
}
