package mysql

import (
	"database/sql"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// func Connect(cfg *MysqlConfig) (*sql.DB, error) {
func Connect(dsn string) (*sql.DB, error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
	//	cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	//)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
