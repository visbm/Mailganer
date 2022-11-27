package handler

import (
	"fmt"
	"html/template"
	"mailganer/internal/models"
	"mailganer/internal/repository"
	"mailganer/pkg/logger"
	"mailganer/pkg/mail"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const datetimeLayout = "2006-01-02T15:04"
const timeLayout = "15:04"

type mailHandler struct {
	logger     logger.Logger
	mail       mail.Mail
	repository *repository.Repository
}

func newMailHandler(logger logger.Logger, mail *mail.Mail , repository *repository.Repository) *mailHandler {
	return &mailHandler{
		logger: logger,
		mail:   *mail,
		repository: repository,
	}
}

var subs = []models.Subscriber{{Address: "samohvalon1998@gmail.com",
	Name: "Nick", Surname: "Fury", FavouriteCategory: "cars"}, {Address: "visbm@mail.ru",
	Name: "kek", Surname: "lol", FavouriteCategory: "films"}}

const (
	home       = "/"
	getSubs    = "/getsubs"
	sendMail   = "/send"
	delaysend  = "/delaysend"
	newsletter = "/newsletter"
)

func (mh *mailHandler) Register(router *mux.Router) {
	router.HandleFunc(home, mh.home).Methods("GET")
	router.HandleFunc(sendMail, mh.sendMail).Methods("POST")
	router.HandleFunc(delaysend, mh.delaysend).Methods("POST")
	router.HandleFunc(newsletter, mh.newsletter).Methods("POST")

}

func (mh *mailHandler) delaysend(w http.ResponseWriter, r *http.Request) {

	/*tmplID, err := strconv.Atoi(r.FormValue("template"))
	if err != nil {
		mh.logger.Errorf("error occurred while getting template. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while getting template. err: %s ", err), http.StatusInternalServerError)
		return
	}*/

	sendTime, err := time.Parse(datetimeLayout, r.FormValue("delay"))
	if err != nil {
		mh.logger.Errorf("error occurred while parsing time. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while parsing time. err: %s ", err), http.StatusBadRequest)
		return
	}

	err = mh.mail.SendMessageWithDelay(sendTime, subs, "templates/mail/hello.html")
	if err != nil {
		mh.logger.Errorf("error occurred while sending message. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while sending message. err: %s ", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("The letter will be sent"))
}

func (mh *mailHandler) newsletter(w http.ResponseWriter, r *http.Request) {
	sendTime, err := time.Parse(timeLayout, r.FormValue("everyday"))
	if err != nil {
		mh.logger.Errorf("error occurred while parsing time. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while parsing time. err: %s ", err), http.StatusBadRequest)
		return
	}

	err = mh.mail.SendMessageEveryDay(sendTime, subs, "templates/mail/hello.html")
	if err != nil {
		mh.logger.Errorf("error occurred while sending message. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while sending message. err: %s ", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Newsletter created"))
}

func (mh *mailHandler) sendMail(w http.ResponseWriter, r *http.Request) {
	var path string
	tmplID, err := strconv.Atoi(r.FormValue("template"))
	if err != nil {
		mh.logger.Errorf("error occurred while getting template. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while getting template. err: %s ", err), http.StatusInternalServerError)
		return
	}

	if tmplID == 1 {
		path = "templates/mail/hello.html"
	} else if tmplID == 2 {
		path = "templates/mail/category.html"

	}

	err = mh.mail.SendMessage(subs, path)
	if err != nil {
		mh.logger.Errorf("error occurred while sending message. err: %s ", err)
		http.Error(w, fmt.Sprintf("error occurred while sending message. err: %s ", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Messages sent"))

}

func (mh *mailHandler) home(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error occurred while parsing template. err: %s ", err), http.StatusInternalServerError)
		mh.logger.Errorf("Can not parse template: %v", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occurred while executing template. err: %s ", err), http.StatusInternalServerError)
		mh.logger.Errorf("Can not execute template: %v", err)
		return
	}
}

func (mh *mailHandler) getSubs(w http.ResponseWriter, r *http.Request) {
	var subs []models.Subscriber

	tmpl, err := template.ParseFiles("templates/subs.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("error occurred while parsing template. err: %s ", err), http.StatusInternalServerError)
		mh.logger.Errorf("Can not parse template: %v", err)
		return
	}

	err = tmpl.Execute(w, subs)
	if err != nil {
		http.Error(w, fmt.Sprintf("error occurred while executing template. err: %s ", err), http.StatusInternalServerError)
		mh.logger.Errorf("Can not execute template: %v", err)
		return
	}
}
