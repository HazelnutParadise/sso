package env

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// todo: 新增.env檔案
var (
	GIN_MODE  string
	LOG_LEVEL logrus.Level

	JWT_SECRET = getEnvOrDefault("JWT_SECRET", "secret key")

	DB_HOST     = getEnvOrDefault("DB_HOST", "localhost")
	DB_PORT     = getEnvOrDefault("DB_PORT", "3306")
	DB_NAME     = getEnvOrDefault("DB_NAME", "sso")
	DB_USER     = getEnvOrDefault("DB_USER", "root")
	DB_PASSWORD = getEnvOrDefault("DB_PASSWORD", "password")
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
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
