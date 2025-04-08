package service

import (
	"errors"
	"snippetbox/internal/modules/user/domain"
	"snippetbox/internal/modules/user/repository"
	"snippetbox/internal/validator"
)

type User struct {
	repo repository.IUserRepository
}

func NewUser(repo repository.IUserRepository) *User {
	return &User{
		repo: repo,
	}
}

func (s *User) UserSignup(form *domain.UsersSignupForm) error {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRGX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		return validator.ErrInvalidForm
	}

	err := s.repo.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")
		}
		return err
	}

	return nil
}

func (s *User) Authenticate(form *domain.UserLoginForm) (int, error) {
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRGX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		return 0, validator.ErrInvalidForm
	}

	id, err := s.repo.Authenticate(form.Email, form.Password)
	if errors.Is(err, repository.ErrInvalidCredentials) {
		form.AddNonFieldError("Email or password is incorrect")
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}
