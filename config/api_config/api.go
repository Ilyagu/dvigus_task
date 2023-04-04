package api_config

import (
	"dvigus_task/config"
	"dvigus_task/internal/pkg/logger"
	"log"
	"os"
	"strconv"
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
		Port: viper.GetString("app.port"),
	}

	RateLimit = config.RateLimitConfig{
		Period:      time.Second * 20,
		MaxRequests: 5,
		Cooldown:    time.Second * 5,
		MaskSize:    24,
	}

	cooldownStr := os.Getenv("COOLDOWN")
	if cooldownStr != "" {
		cl, err := strconv.Atoi(cooldownStr)
		if err == nil {
			RateLimit.Cooldown = time.Second * time.Duration(cl)
		}
	}
	periodStr := os.Getenv("PERIOD")
	if periodStr != "" {
		p, err := strconv.Atoi(periodStr)
		if err == nil {
			RateLimit.Period = time.Second * time.Duration(p)
		}
	}
	maxRequestsStr := os.Getenv("MAX_REQ")
	if maxRequestsStr != "" {
		mxReq, err := strconv.Atoi(maxRequestsStr)
		if err == nil {
			RateLimit.MaxRequests = mxReq
		}
	}
	maskSizeStr := os.Getenv("MASK_SIZE")
	if maskSizeStr != "" {
		m, err := strconv.Atoi(maskSizeStr)
		if err == nil {
			RateLimit.MaskSize = m
		}
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
