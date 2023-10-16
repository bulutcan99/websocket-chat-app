package cmd

import (
	"fmt"
	"github.com/bulutcan99/go-websocket/app/api/controller"
	"github.com/bulutcan99/go-websocket/app/api/middleware"
	"github.com/bulutcan99/go-websocket/app/api/route"
	"github.com/bulutcan99/go-websocket/app/repository"
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
	// Print a message to indicate that the application has started
	fmt.Println("App started")

	// Parse environment variables
	env.ParseEnv()

	// Establish a connection to Redis
	redisClient, err := config.RedisConn()
	if err != nil {
		// Handle connection error and return
		custom_error.ConnectionError()
		return
	}

	// Print a message to indicate a successful Redis connection
	fmt.Println("Redis connected")

	// Defer closing the Redis connection to ensure it's closed when the function exits
	defer redisClient.Close()

	// Establish a connection to the PostgreSQL database
	postgresDB, err := config.PostgresSQLConnection()
	if err != nil {
		// Handle connection error and return
		custom_error.ConnectionError()
		return
	}

	// Print a message to indicate a successful PostgreSQL connection
	fmt.Println("Postgres connected")

	// Defer closing the PostgreSQL database connection
	defer postgresDB.Close()

	// Create an instance of the AuthUserRepo using the PostgreSQL connection
	authRepo := repository.NewAuthUserRepo(postgresDB)

	// Create an AuthController with the AuthUserRepo
	authController := controller.NewAuthController(authRepo)

	// Configure Fiber application settings
	cfg := config.FiberConfig()
	app := fiber.New(cfg)

	// Apply middleware to the Fiber app
	middleware.MiddlewareFiber(app)

	// Serve static files from the "static" directory under "/static" path
	app.Static("/static", "./static")

	// Define routes for the application
	route.Index("/", app)
	route.AuthRoutes(app, authController)

	// Check the value of STAGE_STATUS to determine how to start the server
	if *STAGE_STATUS == stageStatus {
		// Start the server without graceful shutdown
		config.StartServer(app)
	} else {
		// Start the server with graceful shutdown
		config.StartServerWithGracefulShutdown(app)
	}
}
