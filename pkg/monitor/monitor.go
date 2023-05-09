package monitor

import (
	"bcfmonitor/pkg/mail"
	"time"
)

type Monitorizable interface {
	Check() error
	IsUp() bool
	Down() bool
	Up() bool
	Every() time.Duration
}

type Runner struct {
	monitors []*Monitorizable
	mailSvc  *mail.MailService
}

func NewService(mailSvc *mail.MailService) *Runner {
	return &Runner{monitors: make([]*Monitorizable, 0), mailSvc: mailSvc}
}

func (r *Runner) AddMonitorizable(m *Monitorizable) {
	r.monitors = append(r.monitors, m)
}