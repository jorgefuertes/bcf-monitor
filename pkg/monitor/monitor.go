package monitor

import (
	"bcfmonitor/pkg/log"
	"bcfmonitor/pkg/mail"
	"context"
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
	ctx      context.Context
	cancel   context.CancelFunc
	mux      *sync.Mutex
	wg       *sync.WaitGroup
}

func NewService(mailSvc *mail.MailService) *Runner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Runner{
		monitors: []Monitorizable{},
		mailSvc:  mailSvc,
		tickers:  []*time.Ticker{},
		ctx:      ctx,
		cancel:   cancel,
		mux:      &sync.Mutex{},
		wg:       &sync.WaitGroup{},
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
}

func (r *Runner) Run() {
	defer r.Stop()
	for i, m := range r.monitors {
		go r.checkingRoutine(m, r.tickers[i], r.ctx, r.wg)
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
	log.Infof("runner", "Stopping %d monitors", len(r.monitors))
	r.cancel()
	r.wg.Wait()
}

func (r *Runner) checkingRoutine(m Monitorizable, t *time.Ticker, ctx context.Context, wg *sync.WaitGroup) {
	log.Infof("runner", "Starting runner for %s: %s (every %0.2fs)", m.Type(), m.Name(), m.Every().Seconds())
	wg.Add(1)
	for {
		select {
		case <-ctx.Done():
			log.Infof("runner", "Closing runner for %s: %s", m.Type(), m.Name())
			wg.Done()
			return
		case <-t.C:
			err := m.Check()
			if err != nil {
				log.Warnf("runner/err", "%s %s: %s", m.Type(), m.Name(), err)
				if m.IsUp() {
					// was up so I need to alert
					subject := fmt.Sprintf("Outage: Service %s %s is DOWN!", m.Type(), m.Name())
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
				if !m.IsUp() {
					log.Infof("runner", "Recovery: Service %s: %s is UP!", m.Type(), m.Name())
					// was down, notify recover
					subject := fmt.Sprintf("Recovery: Service %s %s is UP again!", m.Type(), m.Name())
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
