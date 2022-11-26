package repository

import (
	"database/sql"
	"mailganer/internal/models"
	"mailganer/pkg/logger"
)

type template struct {
	db     *sql.DB
	logger logger.Logger
}

func newTemplateRepository(db *sql.DB, logger logger.Logger) (repository Template) {
	return &template{
		db:     db,
		logger: logger,
	}
}

func (rep *template) GetTemplateByID(id int) (tmpl *models.Template, err error) {
	tmpl = &models.Template{}
	query := `SELECT 
				template_id, 
				template_path
  			FROM 
			  templates  
  			WHERE 
			  template_id = $1`
	if err = rep.db.QueryRow(query, id).
		Scan(
			&tmpl.ID,
			&tmpl.Path,
		); err != nil {
		rep.logger.Errorf("error occurred while getting template by id, err: %s", err)
		return nil, err
	}

	return tmpl, nil
}
