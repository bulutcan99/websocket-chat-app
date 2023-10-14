package config

import (
	"github.com/bulutcan99/go-websocket/pkg/db/cache"
	db "github.com/bulutcan99/go-websocket/pkg/db/sql"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func StartServerWithGracefulShutdown(a *fiber.App) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := a.Shutdown(); err != nil {
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	fiberConnURL, _ := ConnectionURLBuilder("fiber")
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

func StartServer(a *fiber.App) {
	redisClient, err := cache.RedisConn()
	if err != nil {
		custom_error.ConnectionError()
		return
	}

	defer redisClient.Close()
	postgresDB, err := db.PostgreSQLConnection()
	if err != nil {
		custom_error.ConnectionError()
		return
	}
	defer postgresDB.Close()

	fiberConnURL, _ := ConnectionURLBuilder("fiber")

	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
