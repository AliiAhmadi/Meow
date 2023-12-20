package application

import (
	"Meow/config"
	"Meow/internal/data"
	jlog "Meow/log"
)

type Application struct {
	Config  *config.Config
	Logger  *jlog.Logger
	Version string
	Models  data.Models
}
