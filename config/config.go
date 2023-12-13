package config

import "fmt"

type Config struct {
	Port int
	Env  string
}

func (cfg *Config) GetSport() string {
	return fmt.Sprintf(":%d", cfg.Port)
}
