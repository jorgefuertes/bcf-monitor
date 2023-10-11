package ping

import (
	"bcfmonitor/pkg/monitor/common"
	"fmt"
	"os/exec"
	"runtime"
)

type PingService struct {
	common.MonitorizableBase
	name    string
	host    string
}

func NewService(name, host string, timeout, every int) *PingService {
	s :=  &PingService{name: name, host: host}
	s.TimeoutSeconds = timeout
	s.EverySeconds = every

	return s
}

func (s *PingService) Address() string {
	return s.host
}

func (s *PingService) Check() error {
	timeout := fmt.Sprintf("-i %d", s.TimeoutSeconds)
	if runtime.GOOS == "linux" {
		timeout = fmt.Sprintf("-w %d", s.TimeoutSeconds)
	}
	_, err := exec.Command("ping", s.host, "-c 1", timeout).Output()
	if err != nil {
		s.AddFail()
		return err
	}

	s.Reset()
	return nil
}

func (s *PingService) Type() string {
	return "ping"
}

func (s *PingService) Name() string {
	return s.name
}
