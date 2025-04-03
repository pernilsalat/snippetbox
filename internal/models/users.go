package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	Id             int
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	sqlStatement := `INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(sqlStatement, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) && mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.email") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var (
		id             int
		hashedPassword []byte
	)
	sqlStatement := `SELECT id, hashed_password FROM users WHERE email = ?`
	row := m.DB.QueryRow(sqlStatement, email)
	err := row.Scan(&id, &hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	if err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, ErrInvalidCredentials
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	sqlStatement := `SELECT COUNT(*) FROM users WHERE id = ?`
	err := m.DB.QueryRow(sqlStatement, id).Scan(&exists)

	return exists, err
}
