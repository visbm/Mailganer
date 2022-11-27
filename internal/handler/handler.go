package handler

import (
	"mailganer/internal/repository"
	"mailganer/pkg/logger"
	"mailganer/pkg/mail"

	"github.com/gorilla/mux"
)

type Handler struct {
	logger      logger.Logger
	repository  *repository.Repository
	mailHandler *mailHandler
}

func NewHandler(logger logger.Logger, repository *repository.Repository) *Handler {
	return &Handler{
		logger: logger,
		repository:  repository,
		mailHandler: newMailHandler(logger, mail.NewMail(logger) , repository),
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	h.mailHandler.Register(router)
	return router
}
