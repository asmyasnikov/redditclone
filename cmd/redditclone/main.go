package main

import (
	"github.com/asmyasnikov/redditclone/api"
	"github.com/asmyasnikov/redditclone/internal/app/redditclone"
	"log"
)

func main() {
	cfg, err := api.Get()
	if err != nil {
		log.Println("Can not load the config. Using default config")
		cfg = api.Default()
	}

	if err := redditclone.Run(cfg); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}
