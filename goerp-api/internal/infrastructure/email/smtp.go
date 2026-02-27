package email

import (
	"fmt"
	"net/smtp"
)

type EmailService interface {
	SendCode(to, code string) error
}

type smtpEmailService struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

func NewSMTPService(host string, port int, user, password, from string) EmailService {
	return &smtpEmailService{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
	}
}

func (s *smtpEmailService) SendCode(to, code string) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Login Verification Code\r\n"+
		"\r\n"+
		"Your verification code is: %s. Valid for 5 minutes.\r\n", to, code))

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
