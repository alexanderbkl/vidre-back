package main

import (
	"github.com/alexanderbkl/vidre-back/internal/commands"
	"github.com/alexanderbkl/vidre-back/internal/event"
)
var log = event.Log

func main() {
	commands.Start()
	log.Println("Vidre started.")
}
