package app

import "net/url"

type Pagination struct {
	CurrentPage int
	TotalPages  int
	Pages       []int
	ExtraParams string
}

func NewPagination(currentPage, totalItems, itemsPerPage int, params url.Values) Pagination {
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	// Generate page numbers
	pages := make([]int, 0)
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	// Build extra params string
	params.Del("page")
	extraParams := ""
	if len(params) > 0 {
		extraParams = "&" + params.Encode()
	}

	return Pagination{
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Pages:       pages,
		ExtraParams: extraParams,
	}
}
