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
	"github.com/bulutcan99/go-websocket/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

var (
	Psql          *config.PostgreSQL
	Redis         *config.Redis
	Logger        *zap.Logger
	SchedulerTime time.Duration
	Env           *env.ENV
	stageStatus   = "development"
)

func init() {
	Env = env.ParseEnv()
	Logger = logger.InitLogger(Env.LogLevel)
	Psql = config.NewPostgreSQLConnection()
	zap.S().Info("Postgres connected")
	Redis = config.NewRedisConnection()
	zap.S().Info("Redis connected")

}

func Start() {
	defer Logger.Sync()
	defer Psql.Close()
	defer Redis.Close()
	fmt.Println("App started")

	authRepo := repository.NewAuthUserRepo(Psql)
	redisCache := platform.NewRedisCache(Redis)
	authController := controller.NewAuthController(authRepo, redisCache)
	cfg := config.FiberConfig()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	app.Static("/static", "./static")
	route.Index("/", app)
	route.AuthRoutes(app, authController)
	if Env.StageStatus == stageStatus {
		config.StartServer(app)
	} else {
		config.StartServerWithGracefulShutdown(app)
	}
}
