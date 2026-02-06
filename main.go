package main

import (
	"fmt"
	"log"

	"github.com/db-0/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}
	fmt.Printf("Read configuration file: %+v\n", cfg)

	err = cfg.SetUser("dan")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}

	fmt.Printf("Read configuration file again: %+v\n", cfg)
}
