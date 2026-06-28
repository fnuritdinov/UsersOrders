package models

import (
	errs "UserService/pkg/errors"
)

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
	Role     string
	OTP      string
}

func (r *RegisterRequest) Validate() error {
	if r.Name == "" || r.Email == "" || len(r.Password) < 8 {
		return errs.ErrValidate
	}
	return nil
}

type LoginRequest struct {
	Email    string
	Password string
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" || len(r.Password) < 8 {
		return errs.ErrValidate
	}
	return nil
}
