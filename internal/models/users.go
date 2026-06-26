package models

import (
	errs "UserService/pkg/errors"
	"strings"
	"time"
	"unicode"
)

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

func (u *User) ValidateUser() error {
	if len(u.Name) == 0 && len(u.Email) == 0 && len(u.Password) == 0 && len(u.Role) == 0 {
		return errs.ErrValidate
	}

	return nil
}

func StrEmpty(str string) error {
	if len(str) == 0 {
		return errs.ErrValidate
	}
	return nil
}

func IsValidPhone(phone string) bool {
	if len(phone) != 12 {
		return false
	}

	if !strings.HasPrefix(phone, "992") {
		return false
	}

	for _, ch := range phone {
		if !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}
