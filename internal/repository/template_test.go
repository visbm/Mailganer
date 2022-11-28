package repository

import (
	"errors"
	"mailganer/internal/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_GetTemplateByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := newTemplateRepository(db, log)
	query := `SELECT 
				template_id, 
				template_path
  			FROM 
			  	templates  
  			WHERE 
			  	template_id = $1`

	testTable := []struct {
		name           string
		id             int
		mock           func(id int)
		expectedResult *models.Template
		expectedError  bool
	}{
		{
			name: "OK",
			id:   1,
			mock: func(id int) {
				rows := sqlmock.NewRows([]string{"account_id", "balance"}).
					AddRow(1, "template/tmpl.html")
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)
			},
			expectedResult: &models.Template{
				ID:      1,
				Path: "template/tmpl.html",
			},
			expectedError: false,
		},
		{
			name: "no rows",
			id:   1,
			mock: func(id int) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnError(errors.New("no rows"))
			},

			expectedResult: nil,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.id)
			result, err := r.GetTemplateByID(tt.id)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
