package cmd

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/bulutcan99/go-websocket/app/api/route"
	"github.com/bulutcan99/go-websocket/internal/platform/cache"
	"github.com/bulutcan99/go-websocket/internal/platform/repository"
	config_builder "github.com/bulutcan99/go-websocket/pkg/config"
	config_fiber "github.com/bulutcan99/go-websocket/pkg/config/fiber"
	config_kafka "github.com/bulutcan99/go-websocket/pkg/config/kafka"
	config_psql "github.com/bulutcan99/go-websocket/pkg/config/psql"
	config_redis "github.com/bulutcan99/go-websocket/pkg/config/redis"
	"github.com/bulutcan99/go-websocket/pkg/env"
	"github.com/bulutcan99/go-websocket/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	Psql          *config_psql.PostgreSQL
	Redis         *config_redis.Redis
	KafkaProducer *config_kafka.ProducerKafka
	KafkaConsumer *config_kafka.ConsumerKafka
	Logger        *zap.Logger
	Env           *env.ENV
	stageStatus   = "development"
)

func init() {
	Env = env.ParseEnv()
	Logger = logger.InitLogger(Env.LogLevel)
	Psql = config_psql.NewPostgreSQLConnection()
	Redis = config_redis.NewRedisConnection()
	KafkaProducer = config_kafka.NewKafkaProducerConnection()
	KafkaConsumer = config_kafka.NewKafkaConsumerConnection()
}

// topic         = &env.Env.KafkaMessageTopic
func Start() {
	defer Logger.Sync()
	defer Psql.Close()
	defer Redis.Close()
	defer KafkaProducer.Close()
	defer KafkaConsumer.Close()
	authRepo := repository.NewAuthUserRepo(Psql)
	userRepo := repository.NewUserRepo(Psql)
	redisCache := db_cache.NewRedisCache(Redis)
	// kafkaProducer := pubsub.NewKafkaPublisher(*KafkaProducer)
	// kafkaConsumer := pubsub.NewKafkaSubscriber(*KafkaConsumer)
	authController := controller.NewAuthController(authRepo, redisCache)
	userController := controller.NewUserController(userRepo, redisCache, authController)
	cfg := config_builder.ConfigFiber()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	route.Index("/", app)
	route.AuthRoutes(app, authController)
	route.UserRoutes(app, userController)
	if Env.StageStatus == stageStatus {
		config_fiber.StartServer(app)
	} else {
		config_fiber.StartServerWithGracefulShutdown(app)
	}
}
