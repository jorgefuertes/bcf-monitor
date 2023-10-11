package config

import (
	"bcfmonitor/pkg/mail"
	"os"

	"gopkg.in/yaml.v2"
)

type Database struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	PoolSize int    `json:"poolsize"`
	SSL      bool   `json:"ssl"`
	Timeout  int    `json:"timeout"`
	Every    int    `json:"every"`
	FailAfterRetries int `json:"fail_after_retries" yaml:"fail_after_retries"`
	RememberAfterMinutes int `json:"remember_after_minutes" yaml:"remember_after_retries"`
}

type Redis struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
	Every    int    `json:"every"`
	FailAfterRetries int `json:"fail_after_retries" yaml:"fail_after_retries"`
	RememberAfterMinutes int `json:"remember_after_minutes" yaml:"remember_after_retries"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	FailAfterRetries int `json:"fail_after_retries" yaml:"fail_after_retries"`
	RememberAfterMinutes int `json:"remember_after_minutes" yaml:"remember_after_retries"`
}

type Web struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Needle  string   `json:"needle"`
	Headers []Header `json:"headers"`
	Timeout int      `json:"timeout"`
	Every   int      `json:"every"`
	FailAfterRetries int `json:"fail_after_retries" yaml:"fail_after_retries"`
	RememberAfterMinutes int `json:"remember_after_minutes" yaml:"remember_after_retries"`
}

func (w Web) HeaderMap() map[string]string {
	hm := make(map[string]string)
	for _, h := range w.Headers {
		hm[h.Name] = h.Value
	}

	return hm
}

type SMTP struct {
	Host   string       `json:"host"`
	Port   int          `json:"port"`
	User   string       `json:"user"`
	Pass   string       `json:"pass"`
	Admins []mail.Admin `json:"admins"`
}

type Ping struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Timeout int    `json:"timeout"`
	Every   int    `json:"every"`
	FailAfterRetries int `json:"fail_after_retries" yaml:"fail_after_retries"`
	RememberAfterMinutes int `json:"remember_after_minutes" yaml:"remember_after_retries"`
}

type Configuration struct {
	Databases []Database `json:"databases"`
	Caches    []Redis    `json:"caches"`
	Webs      []Web      `json:"webs"`
	Pings     []Ping     `json:"pings"`
	SMTP      SMTP       `json:"smtp"`
}

func Load(cfgFile string) (cfg Configuration, err error) {
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &cfg)
	return
}

func (cfg Configuration) Dump() (string, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
