package web

import (
	"bcfmonitor/pkg/monitor/common"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WebService struct {
	common.MonitorizableBase
	name    string
	url     string
	needle  string
	headers map[string]string
}

func NewService(name, url, needle string, headers map[string]string, timeout, every int) *WebService {
	s := &WebService{name: name, url: url, needle: needle, headers: headers}
	s.TimeoutSeconds = timeout
	s.EverySeconds = every

	return s
}

func (s *WebService) Address() string {
	return s.url
}

func (s *WebService) Check() error {
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, s.url, nil)
	client.Timeout = s.Timeout()
	if err != nil {
		s.AddFail()
		return err
	}
	for k, v := range s.headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		s.AddFail()
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.AddFail()
		return fmt.Errorf("response: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.AddFail()
		return err
	}
	if !strings.Contains(string(body), s.needle) {
		s.AddFail()
		return fmt.Errorf("body doesn't contains \"%s\"", s.needle)
	}

	s.Reset()
	return nil
}

func (s *WebService) Type() string {
	return "web"
}

func (s *WebService) Name() string {
	return s.name
}
