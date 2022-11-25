package repository

import (
	"database/sql"
	"errors"
	"mailganer/internal/models"
	"mailganer/pkg/logger"
)

var (
	ErrNoRowsAffected = errors.New("now rows affected")
)

type Subscriber interface {
	GetAll() (subs []models.Subscriber, err error)
}

type Repository struct {
	Subscriber
}

func New(db *sql.DB, logger logger.Logger) (repository *Repository) {

	return &Repository{}
}
