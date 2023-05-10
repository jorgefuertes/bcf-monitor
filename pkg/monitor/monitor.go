package monitor

import (
	"bcfmonitor/pkg/log"
	"bcfmonitor/pkg/mail"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Monitorizable interface {
	Type() string
	Name() string
	Address() string
	Check() error
	IsUp() bool
	Down()
	Up()
	Every() time.Duration
}

type Runner struct {
	monitors []Monitorizable
	mailSvc  *mail.MailService
	tickers  []*time.Ticker
	dones    []chan bool
	mux      *sync.Mutex
}

func NewService(mailSvc *mail.MailService) *Runner {
	return &Runner{
		monitors: make([]Monitorizable, 0),
		mailSvc:  mailSvc,
		tickers:  make([]*time.Ticker, 0),
		dones:    make([]chan bool, 0),
		mux:      new(sync.Mutex),
	}
}

func (r *Runner) AddMonitorizable(m Monitorizable) {
	r.mux.Lock()
	defer r.mux.Unlock()
	// setting UP by default
	m.Up()
	// add to runner
	r.monitors = append(r.monitors, m)
	r.tickers = append(r.tickers, time.NewTicker(m.Every()))
	r.dones = append(r.dones, make(chan bool))
}

func (r *Runner) Run() {
	defer r.Stop()
	for i, m := range r.monitors {
		go r.checkingRoutine(m, r.tickers[i], r.dones[i])
	}

	// signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go r.mailSvc.Send("BCFMonitor: Starting", "Monitoring service is starting.\r\n")

	for {
		sig := <-signals
		log.Warnf("runner", "Signal received %s", sig.String())
		return
	}
}

func (r *Runner) Stop() {
	log.Infof("runner", "Stopping %d monitors", len(r.dones))
	for _, done := range r.dones {
		done <- true
	}
}

func (r *Runner) checkingRoutine(m Monitorizable, t *time.Ticker, done chan bool) {
	log.Infof("runner", "Starting runner for %s: %s (every %0.2fs)", m.Type(), m.Name(), m.Every().Seconds())
	for {
		select {
		case <-done:
			log.Infof("runner", "Closing runner for %s: %s", m.Type(), m.Name())
			return
		case <-t.C:
			err := m.Check()
			if err != nil {
				log.Warnf("runner/err", "%s %s: %s", m.Type(), m.Name(), err)
				if m.IsUp() {
					// was up so I need to alert
					subject := fmt.Sprintf("Outage: Service %s %s is down!", m.Type(), m.Name())
					body := fmt.Sprintf("BCFMonitor has detected a service outage:\r\n\r\n"+
						"- Service type: %s\r\n"+
						"- Service name: %s\r\n"+
						"- Service addr: %s\r\n"+
						"\r\n\r\nPlease, take a look here.\r\n",
						m.Type(), m.Name(), m.Address(),
					)
					go r.mailSvc.Send(subject, body)
					m.Down()
				}
			} else {
				log.Infof("runner", "Checking service %s: %s...OK", m.Type(), m.Name())
				if !m.IsUp() {
					// was down, notify recover
					subject := fmt.Sprintf("Recover: Service %s %s is up!", m.Type(), m.Name())
					body := fmt.Sprintf("BCFMonitor has detected a service recovery:\r\n\r\n"+
						"- Service type: %s\r\n"+
						"- Service name: %s\r\n"+
						"- Service addr: %s\r\n"+
						"\r\n\r\nPassed all tests OK.\r\n",
						m.Type(), m.Name(), m.Address(),
					)
					go r.mailSvc.Send(subject, body)
					m.Up()
				}
			}
		}
	}
}
