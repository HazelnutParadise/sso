package env

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// todo: 新增.env檔案
var (
	GIN_MODE   string
	LOG_LEVEL  logrus.Level
	JWT_SECRET string
)

func init() {
	releaseMode := getEnvOrDefault("RELEASE_MODE", "true")
	if releaseMode == "true" {
		GIN_MODE = gin.ReleaseMode
		LOG_LEVEL = logrus.InfoLevel
	} else {
		GIN_MODE = gin.DebugMode
		LOG_LEVEL = logrus.DebugLevel
	}

	JWT_SECRET = getEnvOrDefault("JWT_SECRET", "secret key")
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
