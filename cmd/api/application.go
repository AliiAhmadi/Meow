package main

import (
	"Meow/config"
	"log"
)

type Application struct {
	Config *config.Config
	Logger *log.Logger
}
