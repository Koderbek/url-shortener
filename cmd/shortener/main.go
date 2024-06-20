package main

import (
	"github.com/Koderbek/url-shortener/internal/app/config"
	"github.com/Koderbek/url-shortener/internal/app/server"
)

func main() {
	config.InitConfig()
	server.Start()
}
