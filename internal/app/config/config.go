package config

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var Config struct {
	Flags Flags
}

type Flags struct {
	ServerAddress    string
	ShortenedAddress string
	LogLevel         string
}

func InitConfig() {
	var err error
	Config.Flags, err = initFlags()
	if err != nil {
		panic(err)
	}
}

func initFlags() (Flags, error) {
	var flags Flags
	flag.StringVar(&flags.ServerAddress, "a", ":8080", "address to run server")
	flag.StringVar(&flags.ShortenedAddress, "b", "http://localhost:8080", "base address of the resulting shortened URL")
	flag.StringVar(&flags.LogLevel, "l", "info", "log level")
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		flags.ServerAddress = envServerAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		flags.ShortenedAddress = envBaseURL
	}

	res := strings.Split(flags.ServerAddress, ":")
	if len(res) != 2 {
		return Flags{}, errors.New("invalid server address")
	}

	_, err := strconv.Atoi(res[1])
	if err != nil {
		return Flags{}, errors.New("invalid server address")
	}

	_, err = url.ParseRequestURI(flags.ShortenedAddress)
	if err != nil {
		return Flags{}, errors.New("invalid shortened URL")
	}

	return flags, nil
}
