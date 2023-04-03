package api_config

import (
	"dvigus_task/config"
	"dvigus_task/internal/pkg/logger"
	"log"
	"time"

	"github.com/spf13/viper"
)

var (
	App       config.AppConfig
	RateLimit config.RateLimitConfig
	Logger    logger.Config
)

func SetConfig() {
	viper.SetConfigFile("api_config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	App = config.AppConfig{
		Port:                  viper.GetString("app.port"),
		WellSimilarityPercent: viper.GetInt("app.well_similarity_percent"),
		ExpirationCookieTime:  viper.GetDuration("app.expiration_cookie_time"),
	}

	RateLimit = config.RateLimitConfig{
		Period:      time.Second * 20,
		MaxRequests: 5,
		Cooldown:    time.Second * 5,
		Status:      true,
		MaskSize:    24,
	}

	Logger = logger.Config{
		Path:                viper.GetString("logs.path"),
		Name:                viper.GetString("logs.name"),
		MaxSize:             viper.GetInt("logs.max_size"),
		MaxBackups:          viper.GetInt("logs.max_backups"),
		MaxAge:              viper.GetInt("logs.max_age"),
		RotateCheckInterval: viper.GetDuration("logs.rotate_check_interval"),
	}
}
