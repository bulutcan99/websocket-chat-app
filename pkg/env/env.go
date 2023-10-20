package env

import (
	"fmt"
	custom_error "github.com/bulutcan99/go-websocket/pkg/error"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

type ENV struct {
	StageStatus                    string `env:"STAGE_STATUS,required"`
	ServerHost                     string `env:"SERVER_HOST,required"`
	ServerPort                     int    `env:"SERVER_PORT,required"`
	ServerReadTimeout              int    `env:"SERVER_READ_TIMEOUT,required"`
	JwtSecretKey                   string `env:"JWT_SECRET_KEY,required"`
	JwtSecretKeyExpireMinutesCount int    `env:"JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT,required"`
	JwtRefreshKey                  string `env:"JWT_REFRESH_KEY,required"`
	JwtRefreshKeyExpireHoursCount  int    `env:"JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT,required"`
	DbHost                         string `env:"DB_HOST,required"`
	DbPort                         int    `env:"DB_PORT,required"`
	DbUser                         string `env:"DB_USER,required"`
	DbPassword                     string `env:"DB_PASSWORD,required"`
	DbName                         string `env:"DB_NAME,required"`
	DbSSLMode                      string `env:"DB_SSL_MODE,required"`
	DbMaxConnections               int    `env:"DB_MAX_CONNECTIONS,required"`
	DbMaxIdleConnections           int    `env:"DB_MAX_IDLE_CONNECTIONS,required"`
	DbMaxLifetimeConnections       int    `env:"DB_MAX_LIFETIME_CONNECTIONS,required"`
	RedisHost                      string `env:"REDIS_HOST,required"`
	RedisPort                      int    `env:"REDIS_PORT,required"`
	RedisPassword                  string `env:"REDIS_PASSWORD,required"`
	RedisDBNumber                  int    `env:"REDIS_DB_NUMBER,required"`
	LogLevel                       string `env:"LOG_LEVEL,required"`
}

var doOnce sync.Once
var Env ENV

func ParseEnv() *ENV {
	doOnce.Do(func() {
		e := godotenv.Load()
		if e != nil {
			custom_error.ParseError()
			os.Exit(1)
		}
		if err := env.Parse(&Env); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(0)
		}
	})
	return &Env
}
