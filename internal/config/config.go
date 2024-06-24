package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseURL  string        `yaml:"database_url" env-required:"true"`
	QueryTimeout time.Duration `yaml:"query_timeout" env-default:"2s"`
	Port         string        `yaml:"port" env-required:"true"`
	IsMemoty     bool
}

func Load() (*Config, error) {
	var isMemory bool
	flag.BoolVar(&isMemory, "mem", false, "use memory storage")
	flag.Parse()

	path := "./config/config.yml"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config %s", path)
	}

	cfg.IsMemoty = isMemory
	return &cfg, nil
}
