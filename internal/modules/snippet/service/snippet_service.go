package service

import (
	"snippetbox/internal/modules/snippet/domain"
	"snippetbox/internal/modules/snippet/repository"
	"snippetbox/internal/validator"
)

type Snippet struct {
	repo repository.ISnippetRepository
}

func NewSnippet(repo repository.ISnippetRepository) *Snippet {
	return &Snippet{
		repo: repo,
	}
}

func (s *Snippet) SnippetList() ([]*domain.SnippetModel, error) {
	snippets, err := s.repo.Latest()
	if err != nil {
		return nil, err
	}

	return snippets, nil
}

func (s *Snippet) SnippetRead(id int) (*domain.SnippetModel, error) {
	snippet, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return snippet, nil
}

func (s *Snippet) SnippetCreate(form *domain.SnippetCreateForm) (int, error) {
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field must be between 1 and 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.ValueIn(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		return 0, validator.ErrInvalidForm
	}

	id, err := s.repo.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		return 0, err
	}

	return id, nil
}
