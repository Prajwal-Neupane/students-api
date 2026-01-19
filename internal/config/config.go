package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Server struct {
	Host string
	Address string
}


type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string  `yaml:"storage_path env:"STORAGE_PATH" env-required:"true"`
	Server `yaml:"server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse();

		configPath = *flags


		if configPath == "" {
			log.Fatal("Config is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config

	err:= cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("Cannot read config file %s", err.Error())
	}

	return &cfg
}
