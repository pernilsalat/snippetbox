package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	sqlStatement := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(sqlStatement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	sqlStatement := `SELECT id, title, content, created, expires FROM snippets
		WHERE id = ? AND expires > UTC_TIMESTAMP()`
	row := m.DB.QueryRow(sqlStatement, id)
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	sqlStatement := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return snippets, nil
}
