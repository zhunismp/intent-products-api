package config

type Config struct {
    AppName   string `mapstructure:"APP_NAME"`
    Env       string `mapstructure:"APP_ENV"`
    HTTPPort  int    `mapstructure:"HTTP_PORT"`

    Mongo struct {
        URI      string `mapstructure:"MONGO_URI"`
        Database string `mapstructure:"MONGO_DB"`
    } `mapstructure:",squash"`
}