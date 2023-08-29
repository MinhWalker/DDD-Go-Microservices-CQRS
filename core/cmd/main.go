package main

import (
	"flag"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"log"
)

func main() {
	flag.Parse()

	_, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
}
