package ping

import (
	"bcfmonitor/pkg/log"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

type PingService struct {
	name    string
	host    string
	timeout int
	every   int
	ok      bool
}

func NewService(name, host string, timeout, every int) *PingService {
	return &PingService{name: name, host: host, timeout: timeout, every: every}
}

func (s *PingService) Address() string {
	return s.host
}

func (s *PingService) Check() error {
	timeout := fmt.Sprintf("-i %d", s.timeout)
	if runtime.GOOS == "linux" {
		timeout = fmt.Sprintf("-w %d", s.timeout)
	}
	_, err := exec.Command("ping", s.host, "-c 1", timeout).Output()
	if err != nil {
		return err
	}
	return nil
}

func (s *PingService) IsUp() bool {
	return s.ok
}

func (s *PingService) Down() {
	s.ok = false
	log.Warnf("service/ping", "Service %s is DOWN", s.name)
}

func (s *PingService) Up() {
	s.ok = true
	log.Infof("service/ping", "Service %s is UP", s.name)
}

func (s *PingService) Every() time.Duration {
	return time.Duration(s.every) * time.Second
}

func (s *PingService) Type() string {
	return "ping"
}

func (s *PingService) Name() string {
	return s.name
}
