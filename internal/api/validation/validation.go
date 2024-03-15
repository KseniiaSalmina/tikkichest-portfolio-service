package validation

import (
	"fmt"
	"strconv"

	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/api/response_errors"
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
