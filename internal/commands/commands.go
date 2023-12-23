package commands

import (
	"context"

	"github.com/alexanderbkl/vidre-back/internal/config"
	"github.com/alexanderbkl/vidre-back/internal/event"
	"github.com/alexanderbkl/vidre-back/internal/server"
)

var log = event.Log

func Start() {
	// init logger
	config.InitLogger()

	// load env
	err := config.LoadEnv()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// connect db and define enum types
	err = config.ConnectDB()
	if err != nil {
		log.Fatal("cannot connect to DB and create enums:", err)
	}

	config.InitDb()

	// connect redis
	// config.ConnectRedis()

	// Pass this context down the chain.
	cctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.Start(cctx)
}
