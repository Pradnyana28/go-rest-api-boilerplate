package main

import (
	"flag"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/api"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/config"
)

func main() {

	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
