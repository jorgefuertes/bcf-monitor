package mail

import (
	"fmt"
	"net/smtp"
)

type MailService struct {
	host string
	port int
	user string
	pass string
}

func NewService(host string, port int, user string, pass string) *MailService {
	s := &MailService{host: host, port: port, user: user, pass: pass}
	return s
}

func (s *MailService) Send(to string, subject string, body string) error {
	m := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", s.user, to, subject, body))
	hostAndPort := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.pass, s.host)
	return smtp.SendMail(hostAndPort, auth, s.user, []string{to}, m)
}