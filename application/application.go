package application

import (
	"Meow/config"
	"Meow/internal/data"
	"log"
)

type Application struct {
	Config  *config.Config
	Logger  *log.Logger
	Version string
	Models  data.Models
}
