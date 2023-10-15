package cmd

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/app/api/route"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/bulutcan99/go-websocket/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

var (
	Logger        *zap.Logger
	SchedulerTime time.Duration
	STAGE_STATUS  = &env.Env.StageStatus
	stageStatus   = "development"
)

func Start() {
	fmt.Println("App started")
	env.ParseEnv()
	cfg := config.FiberConfig()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	app.Static("/static", "./static")
	route.Index("/", app)
	route.AuthRoutes(app)
	if *STAGE_STATUS == stageStatus {
		config.StartServer(app)
	} else {
		config.StartServerWithGracefulShutdown(app)
	}
}
