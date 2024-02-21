package response_errors

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"net/http"
)

var ErrMissingID = errors.New("missing id: id is required")
var ErrIncorrectID = errors.New("incorrect id: must be greater than 0")
var ErrIncorrectPortfoliosFilterType = errors.New("incorrect filter: must be ByProfileID, ByCategoryID or empty")

func StatusCodeByErrorWriter(err error, w http.ResponseWriter, ok bool) {
	if errors.Is(err, pgx.ErrNoRows) {
		if !ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)

}
