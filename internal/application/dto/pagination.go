package dto

type Pagination struct {
	CurrentPage int
	PageSize    int
	TotalPages  int
	TotalItems  int
	HasPrevious bool
	HasNext     bool
}

func NewPagination(pageNumber, pageSize, totalPages, totalItems int) *Pagination {
	return &Pagination{
		CurrentPage: pageNumber,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		HasPrevious: pageNumber > 1,
		HasNext:     pageNumber < totalPages,
	}
}
