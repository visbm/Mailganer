package mail

import (
	"errors"
	"fmt"
	"mailganer/pkg/logger"
	"net/smtp"
	"os"
)

// MaillType ...
type MaillType string

// List of Maill Types
const (
	PassReset MaillType = "passReset"
)




// Maill ...
type Maill struct {
	Source   string
	Address  string
	From     string
	Password string
	logger   logger.Logger
}

// GetMaill ...
func GetMaill(source, address, from, password string, logger logger.Logger) *Maill {
	return &Maill{
		Source:   source,
		Address:  address,
		From:     from,
		Password: password,
		logger:   logger,
	}
}

var (
	// ErrInvalidMaill ...
	ErrInvalidMaill = errors.New("invalid mail type")
)

// Auth ...
func (m *Maill) Auth() smtp.Auth {
	return smtp.PlainAuth("", m.From, m.Password, m.Source)
}

// Create ...
func (m *Maill) CreateMessage(linkToReset string, dest []string) (string, error) {
	msg := fmt.Sprint("From: " + m.From + "\n" +
		"To: " + dest[0] + "\n" +
		"Subject: Password reset\n\n" +
		fmt.Sprintf("%s%s \n", os.Getenv("RESTORE_PASSWORD_ENDPOINT"), linkToReset) +
		"Link will be expire in 1 hours.")
	return msg, nil
}

// Send ...
/*func (m *Maill) Send(linkToReset string, dest []string) error {

	msg, err := m.CreateMessage(linkToReset, dest)
	if err != nil {
		m.logger.Errorf("Error occurred while creating message: %v", err)
		return err
	}

	err = smtp.SendMaill(
		m.Address,
		m.Auth(),
		m.From,
		dest,
		[]byte(msg))
	if err != nil {
		m.logger.Errorf("Maill sending error: %v", err)
		return err
	}

	m.logger.Info("Maill send successfull")
	return nil
}
*/