package config

import (
	"fmt"
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
	redisClient, err := RedisConn()
	if err != nil {
		custom_error.ConnectionError()
		return
	}

	fmt.Println("Redis connected")
	defer redisClient.Close()
	postgresDB, err := PostgreSQLConnection()
	if err != nil {
		custom_error.ConnectionError()
		return
	}

	fmt.Println("Postgres connected")
	defer postgresDB.Close()
	fiberConnURL, _ := ConnectionURLBuilder("fiber")
	fmt.Println("Fiber connected")
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
