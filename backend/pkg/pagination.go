package pkg

import "math"

type Pagination struct {
	Page         uint32 `json:"page"`
	PageSize     uint32 `json:"page_size"`
	Total        uint32 `json:"total"`
	TotalPages   uint32 `json:"total_pages"`
	HasNext      bool   `json:"has_next"`
	HasPrevious  bool   `json:"has_previous"`
	NextPage     uint32 `json:"next_page"`
	PreviousPage uint32 `json:"previous_page"`
}

func Offset(page, pageSize uint32) int32 {
	return int32((page - 1) * pageSize)
}

func CalculatePagination(totalItems, pageSize, page uint32) *Pagination {
	if page == 0 {
		page = 1
	}

	totalPages := uint32(math.Ceil(float64(totalItems) / float64(pageSize)))

	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	hasNext := page < totalPages
	hasPrevious := page > 1

	var nextPage uint32
	if hasNext {
		nextPage = page + 1
	} else {
		nextPage = page
	}

	var previousPage uint32
	if hasPrevious {
		previousPage = page - 1
	} else {
		previousPage = 1
	}

	return &Pagination{
		Page:         page,
		PageSize:     pageSize,
		Total:        totalItems,
		TotalPages:   totalPages,
		HasNext:      hasNext,
		HasPrevious:  hasPrevious,
		NextPage:     nextPage,
		PreviousPage: previousPage,
	}
}
