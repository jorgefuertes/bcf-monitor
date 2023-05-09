package main

import (
	"bcfmonitor/pkg/config"
	"bcfmonitor/pkg/log"
	"bcfmonitor/pkg/mail"
	"bcfmonitor/pkg/monitor"
	"bcfmonitor/pkg/monitor/mongo"
	"bcfmonitor/pkg/monitor/redis"
	"bcfmonitor/pkg/monitor/web"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	var CLI struct {
		ConfigFile string `default:"conf/dev.yaml" short:"c" aliases:"conf" help:"path to yaml config file"`
	}
	// command line flags and params
	_ = kong.Parse(&CLI)

	log.Infof("start", "Running on PID: %d", os.Getpid())
	log.Info("config/load", CLI.ConfigFile)
	cfg, err := config.Load(CLI.ConfigFile)
	if err != nil {
		log.Fatalf("config/parse", "ERROR: %s", err)
		os.Exit(1)
	}

	mailSvc := mail.NewService(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.Admins)
	monitorSvc := monitor.NewService(mailSvc)

	// add databases
	for _, d := range cfg.Databases {
		dbMon := mongo.NewService(d.Name, d.Host, d.Port, d.SSL, d.Timeout, d.Every)
		monitorSvc.AddMonitorizable(dbMon)
	}

	// add caches
	for _, c := range cfg.Caches {
		cacheMon := redis.NewService(c.Name, c.Host, c.Port, c.Password, c.Timeout, c.Every)
		monitorSvc.AddMonitorizable(cacheMon)
	}

	// add webs
	for _, w := range cfg.Webs {
		webMon := web.NewService(w.Name, w.URL, w.Needle, w.HeaderMap(), w.Timeout, w.Every)
		monitorSvc.AddMonitorizable(webMon)
	}

	// start the monitor runner
	defer log.Info("runner", "Exiting")
	monitorSvc.Run()
}
