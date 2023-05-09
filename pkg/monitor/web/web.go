package web

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const TIMEOUT_SEC = 5

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
}

func (s *WebService) Up() {
	s.ok = true
}

func (s *WebService) Every() time.Duration {
	return time.Duration(s.every) * time.Second
}