package logger

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func InitLogger(level string) *zap.Logger {
	logLevel := zap.ErrorLevel

	if level == "DEBUG" {
		logLevel = zap.DebugLevel
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)

	return logger
}

func ForceLog(template string, args ...any) {
	template = "%v    INFO    %s " + template + "\n"
	file := getFileAndLine(2)
	args = append([]any{time.Now().Format("2006-01-02T15:04:05.000Z0700"), file}, args...)
	fmt.Printf(template, args...)
}
func getFileAndLine(skip int) string {
	_, file, line, ok := runtime.Caller(skip)

	if !ok {
		return ""
	}

	file = parseFileName(file)

	return file + ":" + strconv.Itoa(line)
}
func parseFileName(file string) string {
	if file == "" {
		return ""
	}

	file = file[len(rootDir()):]

	if file[0] == '/' {
		file = file[1:]
	}

	firstFolderRegex := regexp.MustCompile(`^[a-zA-Z0-9\-]+/`)
	file = firstFolderRegex.ReplaceAllString(file, "$1")

	return file
}
func rootDir() string {
	rootDir, _ := os.Getwd()

	return rootDir
}
