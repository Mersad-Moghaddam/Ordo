package usecase

type PageResult[DataType any] struct {
	Items      []DataType
	TotalCount int64
	PageNumber int
	PageSize   int
}
