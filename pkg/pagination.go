package pkg

import (
	"errors"
	"math"
	"net/http"
	"strconv"
)

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

// ParsePaginationParams extrae y valida page/size de la URL.
// Si no vienen, aplica los valores por defecto recibidos.
func ParsePaginationParams(r *http.Request, defaultPage, defaultSize int) (int, int, error) {
	page := defaultPage
	size := defaultSize

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 15
	}

	if p := r.URL.Query().Get("page"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return 0, 0, errors.New("parametro 'page' invalido")
		}
		page = v
	}

	if s := r.URL.Query().Get("size"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			return 0, 0, errors.New("parametro 'size' invalido")
		}
		size = v
	}

	return page, size, nil
}
