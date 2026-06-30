package errors

import "errors"

var ErrBadRequest = errors.New("bad request")
var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrInternal = errors.New("internal server error")
var ErrValidate = errors.New("error from validate")
