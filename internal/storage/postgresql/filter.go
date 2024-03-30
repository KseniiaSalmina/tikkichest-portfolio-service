package postgresql

import (
	"fmt"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/api/response_errors"
)

type PortfoliosFilterType string

const (
	Empty        = ""
	ByProfileID  = "WHERE portfolios.profile_id = %d"
	ByCategoryID = "WHERE portfolios.category_id = %d"
)

const (
	FilterEmpty        PortfoliosFilterType = ""
	FilterByProfileID  PortfoliosFilterType = "ByProfileID"
	FilterByCategoryID PortfoliosFilterType = "ByCategoryID"
)

var requiredFilters = map[PortfoliosFilterType]string{
	FilterEmpty:        Empty,
	FilterByProfileID:  ByProfileID,
	FilterByCategoryID: ByCategoryID,
}

type PortfoliosFilter struct {
	Type PortfoliosFilterType
	ID   int
}

func portfolioFilter(filterType PortfoliosFilterType, id int) (string, error) {
	filter, ok := requiredFilters[filterType]
	if !ok {
		return "", response_errors.ErrIncorrectPortfoliosFilterType
	}

	switch filterType {
	case FilterEmpty:
		return filter, nil
	default:
		return fmt.Sprintf(filter, id), nil
	}
}
