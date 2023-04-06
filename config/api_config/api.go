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

func SetConfig(filename string) {
	// выставляем кофинг файл
	viper.SetConfigFile(filename)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	// чтение порта из конфиг файла
	App = config.AppConfig{
		Port: viper.GetString("app.port"),
	}

	// выставляем дефолтные значения частоты запросов
	RateLimit = config.RateLimitConfig{
		Period:      time.Second * 20,
		MaxRequests: 5,
		Cooldown:    time.Second * 10,
		MaskSize:    24,
	}

	// если есть данные в переменных окружения, выставляем их
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

	// выставляем конфиг лога из конфиг файла
	Logger = logger.Config{
		Path:                viper.GetString("logs.path"),
		Name:                viper.GetString("logs.name"),
		MaxSize:             viper.GetInt("logs.max_size"),
		MaxBackups:          viper.GetInt("logs.max_backups"),
		MaxAge:              viper.GetInt("logs.max_age"),
		RotateCheckInterval: viper.GetDuration("logs.rotate_check_interval"),
	}
}
