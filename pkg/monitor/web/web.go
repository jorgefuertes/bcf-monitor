package web

import (
	"bcfmonitor/pkg/log"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type WebService struct {
	name    string
	url     string
	needle  string
	timeout int
	every   int
	headers map[string]string
	ok      bool
}

func NewService(name, url, needle string, headers map[string]string, timeout, every int) *WebService {
	return &WebService{name: name, url: url, needle: needle, headers: headers, timeout: timeout, every: every}
}

func (s *WebService) Address() string {
	return s.url
}

func (s *WebService) Check() error {
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, s.url, nil)
	client.Timeout = time.Duration(s.timeout) * time.Second
	if err != nil {
		return err
	}
	for k, v := range s.headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !strings.Contains(string(body), s.needle) {
		return fmt.Errorf("body doesn't contains \"%s\"", s.needle)
	}

	return nil
}

func (s *WebService) IsUp() bool {
	return s.ok
}

func (s *WebService) Down() {
	s.ok = false
	log.Warnf("service/web", "Service %s is DOWN", s.name)
}

func (s *WebService) Up() {
	s.ok = true
	log.Infof("service/web", "Service %s is UP", s.name)
}

func (s *WebService) Every() time.Duration {
	return time.Duration(s.every) * time.Second
}

func (s *WebService) Type() string {
	return "web"
}

func (s *WebService) Name() string {
	return s.name
}
