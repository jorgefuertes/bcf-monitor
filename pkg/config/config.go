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
}

type Redis struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
	Every    int    `json:"every"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Web struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Needle  string   `json:"needle"`
	Headers []Header `json:"headers"`
	Timeout int      `json:"timeout"`
	Every   int      `json:"every"`
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

type Configuration struct {
	Databases []Database `json:"databases"`
	Caches    []Redis    `json:"caches"`
	Webs      []Web      `json:"webs"`
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
