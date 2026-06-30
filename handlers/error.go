package handlers

import (
	errs "UserService/pkg/errors"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

func handleError(w http.ResponseWriter, log *zap.Logger, err error) {
	switch {
	case errors.Is(err, errs.ErrValidate):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, errs.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, errs.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, errs.ErrUnauthorized):
		http.Error(w, err.Error(), http.StatusUnauthorized)

	case errors.Is(err, errs.ErrInternal):
		http.Error(w, err.Error(), http.StatusInternalServerError)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
