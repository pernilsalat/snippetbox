package domain

import (
	"snippetbox/internal/validator"
	"time"
)

type UserModel struct {
	Id             int
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UsersSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
