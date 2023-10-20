package cmd

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/bulutcan99/go-websocket/app/api/route"
	platform "github.com/bulutcan99/go-websocket/db/cache"
	"github.com/bulutcan99/go-websocket/db/repository"
	"github.com/bulutcan99/go-websocket/pkg/config"
	"github.com/bulutcan99/go-websocket/pkg/env"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
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
	redisClient, err := config.RedisConn()
	if err != nil {
		custom_error.ConnectionError()
		return
	}

	fmt.Println("Redis connected")
	defer redisClient.Close()
	postgresDB, err := config.PostgresSQLConnection()
	if err != nil {
		custom_error.ConnectionError()
		return
	}

	fmt.Println("Postgres connected")
	defer postgresDB.Close()
	authRepo := repository.NewAuthUserRepo(postgresDB)
	redisCache := platform.NewRedisCache(redisClient)
	authController := controller.NewAuthController(authRepo, redisCache)
	cfg := config.FiberConfig()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	app.Static("/static", "./static")
	route.Index("/", app)
	route.AuthRoutes(app, authController)
	if *STAGE_STATUS == stageStatus {
		config.StartServer(app)
	} else {
		config.StartServerWithGracefulShutdown(app)
	}
}
