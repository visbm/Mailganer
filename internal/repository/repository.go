package repository

import (
	"database/sql"
	"errors"
	"mailganer/pkg/logger"
)

var (
	ErrNoRowsAffected  = errors.New("now rows affected")
	ErrUserDoesntExist = errors.New("user doesnt exist")
	ErrNewAccNegativeBalance = errors.New("can not set a negative balance for a new account")
)

type Repository struct {
}

func New(db *sql.DB, logger logger.Logger) (repository *Repository) {
	return &Repository{}
}
