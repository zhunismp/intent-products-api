package config

type ServerConfig struct {
	Env           string
	Name          string
	Host          string
	Port          string
	GrpcPort      string
	BaseApiPrefix string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type LoggerConfig struct {
	LogLevel    string
	LogFilePath string
	MaxSize     int
	MaxBackups  int
	MaxAge      int
	Compress    bool
	Endpoint    string
	LogPath     string
}

type AppEnvConfig struct {
	serverCfg *ServerConfig
	dbCfg     *DatabaseConfig
	loggerCfg *LoggerConfig
}
