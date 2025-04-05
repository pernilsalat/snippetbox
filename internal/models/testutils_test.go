package models

import (
	"database/sql"
	_ "embed"
	"testing"
)

//go:embed testdata/setup.sql
var setupDB string

//go:embed testdata/teardown.sql
var teardownDB string

func newTestDB(t *testing.T) (*sql.DB, error) {
	db, err := sql.Open("mysql", "test_web:pass@/test?parseTime=true&multiStatements=true")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	_, err = db.Exec(setupDB)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		_, err := db.Exec(teardownDB)
		if err != nil {
			t.Fatal(err)
		}

		err = db.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	return db, nil
}
