package cmd

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/bulutcan99/go-websocket/app/api/route"
	platform "github.com/bulutcan99/go-websocket/db/cache"
	"github.com/bulutcan99/go-websocket/db/repository"
	"github.com/bulutcan99/go-websocket/pkg/config"
	config_fiber "github.com/bulutcan99/go-websocket/pkg/config/fiber"
	config_psql "github.com/bulutcan99/go-websocket/pkg/config/psql"
	config_redis "github.com/bulutcan99/go-websocket/pkg/config/redis"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/bulutcan99/go-websocket/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	Psql        *config_psql.PostgreSQL
	Redis       *config_redis.Redis
	Logger      *zap.Logger
	Env         *env.ENV
	stageStatus = "development"
)

func init() {
	Env = env.ParseEnv()
	Logger = logger.InitLogger(Env.LogLevel)
	Psql = config_psql.NewPostgreSQLConnection()
	zap.S().Info("Postgres connected")
	Redis = config_redis.NewRedisConnection()
	zap.S().Info("Redis connected")

}

func Start() {
	defer Logger.Sync()
	defer Psql.Close()
	defer Redis.Close()
	zap.S().Info("App started")

	authRepo := repository.NewAuthUserRepo(Psql)
	redisCache := platform.NewRedisCache(Redis)
	authController := controller.NewAuthController(authRepo, redisCache)
	cfg := config.ConfigFiber()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	route.Index("/", app)
	route.AuthRoutes(app, authController)
	if Env.StageStatus == stageStatus {
		config_fiber.StartServer(app)
	} else {
		config_fiber.StartServerWithGracefulShutdown(app)
	}
}
