package response

type PaginatedResponse struct {
	Data       []interface{} `json:"data"`
	Pagination struct {
		CurrentPage int  `json:"current_page"`
		PageSize    int  `json:"page_size"`
		TotalPages  int  `json:"total_pages"`
		TotalItems  int  `json:"total_items"`
		HasPrevious bool `json:"has_previous"`
		HasNext     bool `json:"has_next"`
	} `json:"pagination"`
}
