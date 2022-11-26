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

func newSubscriberRepository(db *sql.DB, logger logger.Logger) (repository Subscriber) {
	return &subscriber{
		db:     db,
		logger: logger,
	}
}

func (rep *subscriber) GetAll() (subs []models.Subscriber, err error) {
	query := `SELECT 	 
				subscriber_id,
				sub_address,
				sub_name,
				sub_surname,
				favourite_category
			FROM 
				subscribers`

	rows, err := rep.db.Query(query)
	if err != nil {
		rep.logger.Errorf("error occurred while getting subscribers. err: %s", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		sub := models.Subscriber{}
		if err = rows.Scan(
			&sub.ID,
			&sub.Address,
			&sub.Name,
			&sub.Surname,
			&sub.FavouriteCategory,
		); err != nil {
			rep.logger.Errorf("error occurred while getting all accounts. err: %s", err)
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}
