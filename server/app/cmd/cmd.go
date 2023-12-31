package cmd

import (
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/bulutcan99/go-websocket/app/api/route"
	"github.com/bulutcan99/go-websocket/internal/platform/cache"
	"github.com/bulutcan99/go-websocket/internal/platform/pubsub"
	repository2 "github.com/bulutcan99/go-websocket/internal/repository"
	wsocket "github.com/bulutcan99/go-websocket/internal/ws"
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

func Start() {
	defer Logger.Sync()
	defer Psql.Close()
	defer Redis.Close()
	defer KafkaProducer.Close()
	defer KafkaConsumer.Close()
	authRepo := repository2.NewAuthUserRepo(Psql)
	userRepo := repository2.NewUserRepo(Psql)
	chatRepo := repository2.NewChatRepo(Psql)
	redisCache := db_cache.NewRedisCache(Redis)
	kafkaProducer := pubsub.NewKafkaPublisher(*KafkaProducer)
	kafkaConsumer := pubsub.NewKafkaSubscriber(*KafkaConsumer)
	newHub := wsocket.NewHub()
	authController := controller.NewAuthController(authRepo, redisCache)
	userController := controller.NewUserController(userRepo, redisCache, authController)
	hubController := controller.NewHubController(newHub, chatRepo, authController, kafkaProducer, kafkaConsumer)
	go newHub.Run()
	cfg := config_builder.ConfigFiber()
	app := fiber.New(cfg)
	middleware.MiddlewareFiber(app)
	route.Index("/", app)
	route.AuthRoutes(app, authController)
	route.UserRoutes(app, userController)
	route.ChatRoutes(app, hubController)
	if Env.StageStatus == stageStatus {
		config_fiber.StartServer(app)
	} else {
		config_fiber.StartServerWithGracefulShutdown(app)
	}
}
