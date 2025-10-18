package config

import (
	"log"

	"github.com/spf13/viper"
)

func Load() *Config {
    v := viper.New()
    v.AutomaticEnv()

    v.SetDefault("APP_NAME", "intent-product-api")
    v.SetDefault("APP_ENV", "development")
    v.SetDefault("HTTP_PORT", 8080)
    v.SetDefault("MONGO_URI", "mongodb://admin:secret@localhost:27017")
    v.SetDefault("MONGO_DB", "product_db")

    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        log.Fatalf("❌ failed to load config: %v", err)
    }

    return &cfg
}