package main

import "ris/internal/app/api"

type Config struct {
	Api api.Config `yaml:"api"`
}

func main() {
	cfg := Config{}
	router := api.New(cfg.Api, nil)
	defer router.Shutdown()

	router.Run()
}
