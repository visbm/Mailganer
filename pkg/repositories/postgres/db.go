package postgres

import (
	"database/sql"
	"fmt"
	"mailganer/pkg/logger"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	Logger   logger.Logger
}

func NewPostgresDB(database *PostgresDB) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		database.Host,
		database.Port,
		database.Username,
		database.DBName,
		database.Password,
		database.SSLMode))
	if err != nil {
		database.Logger.Panicf("Database open error:%s", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		database.Logger.Errorf("DB ping error:%s", err)
		return nil, err
	}
	return db, nil
}
