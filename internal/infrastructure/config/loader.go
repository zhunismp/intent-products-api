package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/zhunismp/intent-products-api/internal/core/infrastructure/config"
)

var _ config.ServerConfigProvider = (*AppEnvConfig)(nil)
var _ config.DatabaseConfigProvider = (*AppEnvConfig)(nil)
var _ config.AppConfigProvider = (*AppEnvConfig)(nil)

var (
	loadedConfig *AppEnvConfig
	loadOnce     sync.Once
)

func LoadConfig(envFilePath ...string) (*AppEnvConfig, error) {
	var loadErr error
	loadOnce.Do(func() {
		if len(envFilePath) > 0 && envFilePath[0] != "" {
			err := godotenv.Load(envFilePath[0])
			if err != nil {
				log.Printf("INFO: .env file not found or failed to load from %s: %v. Proceeding with system environment variables and/or defaults.", envFilePath[0], err)
			}
		}

		serverCfg := &ServerConfig{
			Env:           getEnv("SERVER_ENV", "development"),
			Name:          getEnv("SERVER_NAME", "product-api-dev"),
			Host:          getEnv("SERVER_HOST", "0.0.0.0"),
			Port:          getEnv("SERVER_PORT", "8080"),
			GrpcPort:      getEnv("GRPC_SERVER_PORT", "9000"),
			BaseApiPrefix: getEnv("SERVER_BASEAPIPREFIX", "/api/v1"),
		}

		dbCfg := &DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "27017"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "secret"),
			Name:     getEnv("DB_NAME", "product_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Bangkok"),
		}

		// Parse logger config values
		maxSize := mustParseInt(getEnv("LOGGING_MAXSIZE", "100"), "LOGGING_MAXSIZE")
		maxBackups := mustParseInt(getEnv("LOGGING_MAXBACKUPS", "3"), "LOGGING_MAXBACKUPS")
		maxAge := mustParseInt(getEnv("LOGGING_MAXAGE", "30"), "LOGGING_MAXAGE")
		compress := mustParseBool(getEnv("LOGGING_COMPRESS", "true"), "LOGGING_COMPRESS")

		loggerCfg := &LoggerConfig{
			LogLevel:    getEnv("LOGGING_LEVEL", "INFO"),
			LogFilePath: getEnv("LOGGING_PATH", "/var/log/app.log"),
			MaxSize:     maxSize,
			MaxBackups:  maxBackups,
			MaxAge:      maxAge,
			Compress:    compress,
			Endpoint:    getEnv("LOGGING_ENDPOINT", "0.0.0.0:3100/otlp/v1/logs"),
		}

		if dbCfg.User == "" {
			loadErr = fmt.Errorf("DB_USER cannot be empty")
			return
		}

		loadedConfig = &AppEnvConfig{
			serverCfg: serverCfg,
			dbCfg:     dbCfg,
			loggerCfg: loggerCfg,
		}

		log.Println("INFO: Application configuration loaded successfully.")
	})

	if loadErr != nil {
		return nil, loadErr
	}
	if loadedConfig == nil && loadErr == nil {
		return nil, fmt.Errorf("configuration was not loaded")
	}
	return loadedConfig, nil
}

func mustParseInt(val, key string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return i
}

func mustParseBool(val, key string) bool {
	b, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return b
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

/* Application Cfg */
func (c *AppEnvConfig) GetServerEnv() string           { return c.serverCfg.Env }
func (c *AppEnvConfig) GetServerName() string          { return c.serverCfg.Name }
func (c *AppEnvConfig) GetServerHost() string          { return c.serverCfg.Host }
func (c *AppEnvConfig) GetServerPort() string          { return c.serverCfg.Port }
func (c *AppEnvConfig) GetServerBaseApiPrefix() string { return c.serverCfg.BaseApiPrefix }
func (c *AppEnvConfig) GetGrpcServerPort() string      { return c.serverCfg.GrpcPort }

/* Database Cfg */
func (c *AppEnvConfig) GetDBHost() string     { return c.dbCfg.Host }
func (c *AppEnvConfig) GetDBPort() string     { return c.dbCfg.Port }
func (c *AppEnvConfig) GetDBUser() string     { return c.dbCfg.User }
func (c *AppEnvConfig) GetDBPassword() string { return c.dbCfg.Password }
func (c *AppEnvConfig) GetDBName() string     { return c.dbCfg.Name }
func (c *AppEnvConfig) GetDBSSLMode() string  { return c.dbCfg.SSLMode }
func (c *AppEnvConfig) GetDBTimezone() string { return c.dbCfg.Timezone }
func (c *AppEnvConfig) GetDBDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.dbCfg.Host, c.dbCfg.Port, c.dbCfg.User, c.dbCfg.Password, c.dbCfg.Name, c.dbCfg.SSLMode, c.dbCfg.Timezone)
}

/* Logging Cfg */
func (c *AppEnvConfig) GetLogLevel() string    { return c.loggerCfg.LogLevel }
func (c *AppEnvConfig) GetLogFilePath() string { return c.loggerCfg.LogFilePath }
func (c *AppEnvConfig) GetMaxSize() int        { return c.loggerCfg.MaxSize }
func (c *AppEnvConfig) GetMaxBackups() int     { return c.loggerCfg.MaxBackups }
func (c *AppEnvConfig) GetMaxAge() int         { return c.loggerCfg.MaxAge }
func (c *AppEnvConfig) GetCompress() bool      { return c.loggerCfg.Compress }
func (c *AppEnvConfig) GetLogEndpoint() string { return c.loggerCfg.Endpoint }
