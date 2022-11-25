package mail

import (
	"errors"
	"html/template"
	"io"
	"mailganer/internal/models"
	"mailganer/pkg/logger"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	host   string
	port   string
	name   string
	pass   string
	logger logger.Logger
}

func NewMail(logger logger.Logger) *Mail {
	return &Mail{
		host:   os.Getenv("MAIL_SOURCE"),
		port:   os.Getenv("MAIL_PORT"),
		name:   os.Getenv("MAIL_FROM"),
		pass:   os.Getenv("MAIL_PASSWORD"),
		logger: logger,
	}
}

func (m *Mail) SendMessage(sub []models.Subscriber) error {
	port, err := strconv.Atoi(m.port)
	if err != nil {
		m.logger.Errorf("error occurred while conv port. err:%s ", err)
		return err
	}

	d := gomail.NewDialer(m.host, port, m.name, m.pass)
	s, err := d.Dial()
	if err != nil {
		m.logger.Errorf("error occurred while dial to smtp server. err:%s ", err)
		return err
	}
	defer s.Close()

	mess := gomail.NewMessage()

	tmpl, err := template.ParseFiles("templates/hello.html")
	if err != nil {
		m.logger.Errorf("error occurred while parse template. err:%s ", err)
		return err
	}
	for _, r := range sub {
		mess.SetHeader("From", m.name)
		mess.SetAddressHeader("To", r.Address, r.Name)
		mess.SetHeader("Subject", "MailGaner")
		mess.AddAlternativeWriter("text/html", func(w io.Writer) error {
			return tmpl.Execute(w, r)
		})
		if err := gomail.Send(s, mess); err != nil {
			m.logger.Errorf("Could not send email to %q: %v", r.Address, err)
		}
		mess.Reset()
	}
	return nil

}

func (m *Mail) SendMessageWithDelay(callTime time.Time, sub []models.Subscriber) error {
	
	now := time.Now()
	if callTime.Before(now) {
		return errors.New("time can not be befor now")
	}

	duration := callTime.Sub(now)

	go func() {
		time.Sleep(duration)
		m.SendMessage(sub)
	}()

	return nil
}

func (m *Mail) SendMessageEveryDay(callTime time.Time, sub []models.Subscriber) error {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		m.logger.Errorf("error occurred while load location. err:%s ", err)
		return err
	}
	now := time.Now().Local()
	firstCallTime := time.Date(
		now.Year(), now.Month(), now.Day(), callTime.Hour(), callTime.Minute(), 0, 0, loc)
	if firstCallTime.Before(now) {
		firstCallTime = firstCallTime.Add(time.Hour * 24)
	}

	duration := firstCallTime.Sub(time.Now().Local())

	go func() {
		time.Sleep(duration)
		for {
			m.SendMessage(sub)
			time.Sleep(time.Hour * 24)
		}
	}()

	return nil
}
