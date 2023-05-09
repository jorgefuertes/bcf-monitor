package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Database struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	PoolSize int    `json:"poolsize"`
	SSL      bool   `json:"ssl"`
}

type Redis struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type Web struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Needle string `json:"needle"`
}

type Admin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SMTP struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Configuration struct {
	Databases    []Database
	Caches       []Redis
	Applications []Web
	SMTP         SMTP `json:"smtp"`
	Admins       []Admin
}

func Load(cfgFile string) (cfg Configuration, err error) {
	content, err := os.ReadFile(cfgFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &cfg)
	return
}