package config_fiber

import (
	"github.com/bulutcan99/go-websocket/pkg/config"
	"go.uber.org/zap"
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

	fiberConnURL, _ := config_builder.ConnectionURLBuilder("fiber")
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

func StartServer(a *fiber.App) {
	fiberConnURL, _ := config_builder.ConnectionURLBuilder("fiber")
	zap.S().Info("Connected to Fiber successfully.")
	if err := a.Listen(fiberConnURL); err != nil {
		zap.S().Errorf("Oops... Server is not running! Reason: %v", err)
	}
}
