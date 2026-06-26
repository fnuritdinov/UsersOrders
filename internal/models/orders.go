package models

import (
	errs "UserService/pkg/errors"
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Product   string    `json:"product"`
	Price     int       `json:"price"`
	UserID    int       `json:"userID"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (o *Order) ValidateCreateOrder() error {
	if len(o.Product) == 0 {
		return errs.ErrValidate
	}
	if o.Price < 1 {
		return errs.ErrValidate
	}
	if o.UserID < 1 {
		return errs.ErrValidate
	}
	return nil
}
