package service

import (
	"snippetbox/internal/user/domain"
	"snippetbox/internal/user/repository"
	"snippetbox/internal/validator"
)

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) UserSignup(form *domain.UsersSignupForm) (bool, error) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRGX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		return false, nil
	}

	err := s.repo.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		return false, err
	}

	return true, nil
}
