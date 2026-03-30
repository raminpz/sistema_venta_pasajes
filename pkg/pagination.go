package pkg

import "math"

type PaginationMeta struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Paginate calcula offset y limit para paginación y retorna metadata
func Paginate(page, size, total int) (offset, limit int, meta PaginationMeta) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 15
	}
	offset = (page - 1) * size
	limit = size
	totalPages := int(math.Ceil(float64(total) / float64(size)))
	meta = PaginationMeta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}
	return
}
