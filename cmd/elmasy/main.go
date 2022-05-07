package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/elmasy-com/elmasy/internal/config"
	"github.com/elmasy-com/elmasy/internal/router"
)

func main() {

	confPath := flag.String("config", "./elmasy.conf", "Path to the config file")
	flag.Parse()

	if err := config.ParseConfig(*confPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %s\n", err)
		os.Exit(1)
	}

	r := router.SetupRouter()

	if config.GlobalConfig.SSLCertificate == "" || config.GlobalConfig.SSLKey == "" {
		if err := r.Run(config.GlobalConfig.Listen); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run server: %s\n", err)
		}
	} else {
		err := r.RunTLS(config.GlobalConfig.Listen, config.GlobalConfig.SSLCertificate, config.GlobalConfig.SSLKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to run server: %s\n", err)
		}
	}
}
