package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/elmasy-com/elmasy/config"
)

func main() {

	confPath := flag.String("config", "./elmasy.conf", "Path to the config file")
	flag.Parse()

	if err := config.ParseConfig(*confPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %s\n", err)
		os.Exit(1)
	}

	r := SetupRouter()

	if err := r.Run(config.GlobalConfig.ListenAddr); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run server: %s\n", err)
	}
}
