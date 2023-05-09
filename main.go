package main

import (
	"bcfmonitor/pkg/config"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	var CLI struct {
		ConfigFile string `default:"conf/dev.yaml" short:"c" aliases:"conf" help:"path to yaml config file"`
	}
	// command line flags and params
	_ = kong.Parse(&CLI)

	_, err := config.Load(CLI.ConfigFile)
	if err != nil {
		fmt.Printf("[CONFIG/PARSE] %s\n", CLI.ConfigFile)
		fmt.Printf("[CONFIG/PARSE] %s\n", err)
		os.Exit(1)
	}

	fmt.Println("OK")
}