package email

import (
	"fmt"
	"github.com/pkg/errors"
	"net/smtp"
	"strings"
)

type SMTP struct {
	Server   string		`yaml:"server"`
	Port     string   	`yaml:"port"`
	Password string   	`yaml:"password"`
	Username string   	`yaml:"username"`
	From     string   	`yaml:"from"`
	To       []string 	`yaml:"to"`
	Subject  string		`yaml:"subject"`
}

func Send(subject string, body string, config *SMTP) error {
	msg := "From: " + config.From + "\r\n" +
		"To: " + strings.Join(config.To, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n"

	addr := fmt.Sprintf("%v:%v", config.Server, config.Port)
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Server)

	if err := smtp.SendMail(addr, auth, config.From, config.To, []byte(msg)); err != nil {
		return errors.Wrapf(err, "sending email notification failed")
	}
	return nil
}
