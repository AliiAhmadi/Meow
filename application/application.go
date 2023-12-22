package application

import (
	"Meow/config"
	"Meow/internal/data"
	jlog "Meow/log"
	"Meow/mailer"
	"sync"
)

type Application struct {
	Config  *config.Config
	Logger  *jlog.Logger
	Version string
	Models  data.Models
	Mailer  mailer.Mailer
	Wg      sync.WaitGroup
}
