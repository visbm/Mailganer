package repository

import (
	"database/sql"
	"mailganer/internal/models"
	"mailganer/pkg/logger"
)

type subscriber struct {
	db     *sql.DB
	logger logger.Logger
}

func NewAccountRepository(db *sql.DB, logger logger.Logger) (repository Subscriber) {
	return &subscriber{
		db:     db,
		logger: logger,
	}
}

func (rep *subscriber) GetAll() (ubs []models.Subscriber, err error) {
	query := `SELECT id,	 
				adress,
				name,
				surname,
				favourite_category
			FROM subscribers`

	rows, err := rep.db.Query(query)
	if err != nil {
		rep.logger.Errorf("error occurred while getting subscribers. err: %s", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		sub := models.Subscriber{}
		if err = rows.Scan(
			&account.ID,
			&account.Balance,
		); err != nil {
			rep.logger.Errorf("error occurred while getting all accounts. err: %s", err)
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
