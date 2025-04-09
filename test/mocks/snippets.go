package mocks

import (
	"snippetbox/internal/modules/snippet/domain"
	"snippetbox/internal/platform/database"
	"time"
)

var mockSnippet = &domain.SnippetModel{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*domain.SnippetModel, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, database.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*domain.SnippetModel, error) {
	return []*domain.SnippetModel{mockSnippet}, nil
}
