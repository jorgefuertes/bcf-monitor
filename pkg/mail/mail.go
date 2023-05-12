package mail

import (
	"bcfmonitor/pkg/log"
	"fmt"
	"net/smtp"
)

type MailService struct {
	host   string
	port   int
	user   string
	pass   string
	admins []Admin
}

type Admin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewService(host string, port int, user string, pass string, admins []Admin) *MailService {
	s := &MailService{host: host, port: port, user: user, pass: pass, admins: admins}
	return s
}

func (s *MailService) Send(subject string, body string) {
	log.Infof("mail/send", "Sending alert to admins: \"%s\"", subject)
	for _, a := range s.admins {
		m := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", s.user, a.Email, subject, body))
		hostAndPort := fmt.Sprintf("%s:%d", s.host, s.port)
		auth := smtp.PlainAuth("", s.user, s.pass, s.host)
		err := smtp.SendMail(hostAndPort, auth, s.user, []string{a.Email}, m)
		if err != nil {
			log.Error("mail/send", err.Error())
		}
	}
}
