package repository

import (
	"errors"
	"mailganer/internal/models"
	"mailganer/pkg/logger"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var log = logger.GetLogger()

func Test_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := newSubscriberRepository(db, log)

	query := `SELECT 	 
				subscriber_id,
				sub_address,
				sub_name,
				sub_surname,
				favourite_category
			FROM 
				subscribers`

	testTable := []struct {
		name           string
		mock           func()
		expectedResult []models.Subscriber
		expectedError  bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"subscriber_id", "sub_address", "sub_name", "sub_surname", "favourite_category"}).
					AddRow(1, "email@gmail.com", "John", "Smith", "cars").AddRow(2, "email2@gmail.com", "Jack", "Sparrow", "ships")
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
			},
			expectedResult: []models.Subscriber{
				{
					ID:                1,
					Address:           "email@gmail.com",
					Name:              "John",
					Surname:           "Smith",
					FavouriteCategory: "cars",
				},
				{
					ID:                2,
					Address:           "email2@gmail.com",
					Name:              "Jack",
					Surname:           "Sparrow",
					FavouriteCategory: "ships",
				},
			},

			expectedError: false,
		},
		{
			name: "no rows",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(errors.New("no rows"))
			},

			expectedResult: nil,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			result, err := r.GetAll()
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
