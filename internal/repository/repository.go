package repository

import (
	"database/sql"
	"mailganer/internal/models"
	"mailganer/pkg/logger"
	 _ "github.com/golang/mock/mockgen/model"
)

//go:generate mockgen -source repository.go -destination mock_repository/mock_repository.go 

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
