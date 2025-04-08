package repository

import (
	"database/sql"
	"errors"
	"snippetbox/internal/modules/snippet/domain"
	"snippetbox/internal/platform/database"
)

type ISnippetRepository interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (*domain.SnippetModel, error)
	Latest() ([]*domain.SnippetModel, error)
}

type SnippetRepository struct {
	DB database.DB
}

func NewSnippet(db database.DB) *SnippetRepository {
	return &SnippetRepository{
		DB: db,
	}
}

func (r *SnippetRepository) Insert(title, content string, expires int) (int, error) {
	sqlStatement := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := r.DB.Exec(sqlStatement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *SnippetRepository) Get(id int) (*domain.SnippetModel, error) {
	sqlStatement := `SELECT id, title, content, created, expires FROM snippets
		WHERE id = ? AND expires > UTC_TIMESTAMP()`
	row := r.DB.QueryRow(sqlStatement, id)
	s := &domain.SnippetModel{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *SnippetRepository) Latest() ([]*domain.SnippetModel, error) {
	sqlStatement := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := r.DB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*domain.SnippetModel
	for rows.Next() {
		s := &domain.SnippetModel{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	return snippets, nil
}
