package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/honyshyota/constanta-rest-api/internal/app/apiserver"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	// Parsing conf from command line
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	
	if err := apiserver.Start(config); err != nil {
		logrus.Fatal("connection db error: ", err)
	}
}
