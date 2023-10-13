package cmd

import (
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/bulutcan99/go-websocket/app/api/route"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

var Logger *zap.Logger
var SchedulerTime time.Duration
var STAGE_STATUS = &env.Env.StageStatus

func InitLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	return logger
}

func Init() {
	Logger = InitLogger()
	_ = env.ParseEnv()
}

func Start() {
	Init()
	defer func(Logger *zap.Logger) {
		if r := recover(); r != nil {
			Logger.Error("Panic occurred:", zap.Any("error", r))
		}
		err := Logger.Sync()
		if err != nil {
			zap.S().Debug("There is an error while logging!")
		}
	}(Logger)

	cfg := config.FiberConfig()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	app.Static("/static", "./static")
	route.Index("/", app)

	if *STAGE_STATUS == "development" {
		config.StartServer(app)
	} else {
		config.StartServerWithGracefulShutdown(app)
	}
}
