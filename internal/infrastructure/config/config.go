package infrastructure

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

type AppEnvConfig struct {
	serverCfg *ServerConfig
	dbCfg     *DatabaseConfig
}
