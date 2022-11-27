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

type Template interface {
	GetTemplateByID(id int) (acc *models.Template, err error)
}

type Repository struct {
	Subscriber
	Template
}

func New(db *sql.DB, logger logger.Logger) (repository *Repository) {
	return &Repository{
		Subscriber: newSubscriberRepository(db, logger),
		Template:   newTemplateRepository(db, logger),
	}
}
