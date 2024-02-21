package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type pageInfo struct {
	number int
	limit  int
	offset int
}

const defaultLimit int = 30

func (s *Server) getPageInfo(r *http.Request) (*pageInfo, error) {
	limitStr := r.FormValue("limit")
	var limit int
	switch limitStr {
	case "":
		limit = defaultLimit
	default:
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("failed to get limit: %w", err)
		}
		limit = l
	}
	if limit <= 0 {
		return nil, errors.New("limit must be greater than 0")
	}

	pageNoStr := r.FormValue("page")
	var page int
	switch pageNoStr {
	case "":
		page = 1
	default:
		p, err := strconv.Atoi(pageNoStr)
		if err != nil {
			return nil, fmt.Errorf("failed to get page number: %w", err)
		}
		page = p
	}
	if page <= 0 {
		return nil, errors.New("page number must be greater than 0")
	}

	result := pageInfo{number: page, limit: limit, offset: (page - 1) * limit}
	return &result, nil
}
