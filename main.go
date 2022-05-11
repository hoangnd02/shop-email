package main

import (
	"github.com/hoanggggg5/shopemail/config"
	"github.com/hoanggggg5/shopemail/services"
)

func main() {
	config.InnitConfig()

	services.NewSendEmail().Process()
}
