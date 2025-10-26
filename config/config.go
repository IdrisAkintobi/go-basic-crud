package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Port                 string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	SessionDuration      int
	SessionRefreshWindow int
	TokenLength          int
	MaximumSession       int
	GeoIPAccountID       string
	GeoIPLicenseKey      string
	TestDatabaseName     string
	TestDatabaseURL      string
	DataFilePath         string
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		sessionDuration, _ := strconv.Atoi(getEnvOrDefault("SESSION_DURATION", "60"))
		sessionRefreshWindow, _ := strconv.Atoi(getEnvOrDefault("SESSION_REFRESH_WINDOW", "10"))
		tokenLength, _ := strconv.Atoi(getEnvOrDefault("TOKEN_LENGTH", "32"))
		maxSession, _ := strconv.Atoi(getEnvOrDefault("MAXIMUM_SESSION", "5"))

		instance = &Config{
			Port:                 getEnvOrDefault("PORT", "3003"),
			DBHost:               os.Getenv("DB_HOST"),
			DBPort:               os.Getenv("DB_PORT"),
			DBUser:               os.Getenv("DB_USER"),
			DBPassword:           os.Getenv("DB_PASSWORD"),
			DBName:               os.Getenv("DB_NAME"),
			SessionDuration:      sessionDuration,
			SessionRefreshWindow: sessionRefreshWindow,
			TokenLength:          tokenLength,
			MaximumSession:       maxSession,
			GeoIPAccountID:       os.Getenv("GEO21P_ACCOUNT_ID"),
			GeoIPLicenseKey:      os.Getenv("GEO21P_LICENSE_KEY"),
			TestDatabaseName:     os.Getenv("TEST_DATABASE_NAME"),
			TestDatabaseURL:      os.Getenv("TEST_DATABASE_URL"),
			DataFilePath:         getEnvOrDefault("DATA_FILE_PATH", "./geo2ip-data/GeoLite2-City.mmdb"),
		}
	})
	return instance
}

func Get() *Config {
	if instance == nil {
		return Load()
	}
	return instance
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
