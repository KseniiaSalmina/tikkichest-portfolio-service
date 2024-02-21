package validation

import (
	"fmt"
	"strconv"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/api/response_errors"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/storage/postgresql"
)

func ID(idStr string) (int, error) {
	if idStr == "" {
		return 0, response_errors.ErrMissingID
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("incorrect id format: %w", err)
	}
	if id <= 0 {
		return 0, response_errors.ErrIncorrectID
	}

	return id, nil
}

func PortfoliosFilter(filterType string, id int) (*postgresql.PortfoliosFilter, error) {
	var filter postgresql.PortfoliosFilter
	switch filterType {
	case "":
	case "ByProfileID":
		filter = postgresql.PortfoliosFilter{ID: id, Type: postgresql.ByProfileID}
	case "ByCategoryID":
		filter = postgresql.PortfoliosFilter{ID: id, Type: postgresql.ByCategoryID}
	default:
		return nil, response_errors.ErrIncorrectPortfoliosFilterType
	}

	return &filter, nil
}
