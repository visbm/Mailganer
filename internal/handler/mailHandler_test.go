package handler

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"mailganer/internal/models"
	"mailganer/internal/repository/mock_repository"
	"mailganer/pkg/logger"
	"mailganer/pkg/mail"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var log = logger.GetLogger()

func Test_getSubs(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockSubscriber)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "error parse template",
			mockBehavior: func(s *mock_repository.MockSubscriber) {
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while parsing template. err: open templates/subs.html: The system cannot find the path specified. \n",
		},
		{
			name: "no accounts",

			mockBehavior: func(s *mock_repository.MockSubscriber) {
				s.EXPECT().GetAll().Return(nil, errors.New("sql: no rows in result set"))

			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting subscribers. err: sql: no rows in result set \n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			subRep := mock_repository.NewMockSubscriber(c)
			tmplRep := mock_repository.NewMockTemplate(c)

			testCase.mockBehavior(subRep)
			mail := mail.NewMail(log)
			handler := newMailHandler(log, mail, subRep, tmplRep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("GET", "/getsubs", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_Home(t *testing.T) {
	testTable := []struct {
		name                string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "error parse template",

			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while parsing template. err: open templates/home.html: The system cannot find the path specified. \n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			subRep := mock_repository.NewMockSubscriber(c)
			tmplRep := mock_repository.NewMockTemplate(c)

			mail := mail.NewMail(log)
			handler := newMailHandler(log, mail, subRep, tmplRep)
			router := mux.NewRouter()
			handler.Register(router)

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_sendMail(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "error parse template",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while sending message. err: strconv.Atoi: parsing \"\": invalid syntax \n",
		},
		{
			name: "error sending message",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while sending message. err: strconv.Atoi: parsing \"\": invalid syntax \n",
		},
		{
			name: "no accounts",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return(nil, errors.New("sql: no rows in result set"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting subscribers. err: sql: no rows in result set \n",
		},
		{
			name: "no template",

			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(nil, errors.New("sql: no rows in result set"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting template. err: sql: no rows in result set \n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			subRep := mock_repository.NewMockSubscriber(c)
			tmplRep := mock_repository.NewMockTemplate(c)

			testCase.mockBehavior(subRep, tmplRep)
			mail := mail.NewMail(log)
			handler := newMailHandler(log, mail, subRep, tmplRep)
			router := mux.NewRouter()
			handler.Register(router)

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			writer.WriteField("template", "1")
			writer.Close()

			req := httptest.NewRequest("POST", "/send", payload)
			w := httptest.NewRecorder()
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func Test_delaysend(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate)

	testTable := []struct {
		name, time          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "error sending message",
			time: "2006-01-02T15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while sending message. err: time can not be befor now \n",
		},
		{
			name: "ok",
			time: "2023-01-02T15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "The letter will be sent",
		},
		{
			name: "no accounts",
			time: "2023-01-02T15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return(nil, errors.New("sql: no rows in result set"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting subscribers. err: sql: no rows in result set \n",
		},
		{
			name: "no template",
			time: "2023-01-02T15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(nil, errors.New("sql: no rows in result set"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting template. err: sql: no rows in result set \n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			subRep := mock_repository.NewMockSubscriber(c)
			tmplRep := mock_repository.NewMockTemplate(c)

			testCase.mockBehavior(subRep, tmplRep)
			mail := mail.NewMail(log)
			handler := newMailHandler(log, mail, subRep, tmplRep)
			router := mux.NewRouter()
			handler.Register(router)

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			writer.WriteField("template", "1")
			writer.WriteField("delay", testCase.time)
			writer.Close()

			req := httptest.NewRequest("POST", "/delaysend", payload)
			w := httptest.NewRecorder()
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}


func Test_newsletter(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate)

	testTable := []struct {
		name, time          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "ok",
			time: "15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "Newsletter created",
		},
		{
			name: "empty form",
			time: "",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return([]models.Subscriber{{}}, nil)
			},
			expectedStatusCode:  400,
			expectedRequestBody: "error occurred while parsing time. err: parsing time \"\" as \"15:04\": cannot parse \"\" as \"15\" \n",
		},
		{
			name: "no accounts",
			time: "15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(&models.Template{}, nil)
				s.EXPECT().GetAll().Return(nil, errors.New("sql: no rows in result set"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting subscribers. err: sql: no rows in result set \n",
		},
		{
			name: "no template",
			time: "15:04",
			mockBehavior: func(s *mock_repository.MockSubscriber, t *mock_repository.MockTemplate) {
				t.EXPECT().GetTemplateByID(1).Return(nil, errors.New("sql: no rows in result set"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "error occurred while getting template. err: sql: no rows in result set \n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			subRep := mock_repository.NewMockSubscriber(c)
			tmplRep := mock_repository.NewMockTemplate(c)

			testCase.mockBehavior(subRep, tmplRep)
			mail := mail.NewMail(log)
			handler := newMailHandler(log, mail, subRep, tmplRep)
			router := mux.NewRouter()
			handler.Register(router)

			payload := &bytes.Buffer{}
			writer := multipart.NewWriter(payload)
			writer.WriteField("template", "1")
			writer.WriteField("everyday", testCase.time)
			writer.Close()

			req := httptest.NewRequest("POST", "/newsletter", payload)
			w := httptest.NewRecorder()
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}