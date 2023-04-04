package config

import (
	"time"
)

type AppConfig struct {
	Port string
}

type RateLimitConfig struct {
	Period      time.Duration
	MaxRequests int
	Cooldown    time.Duration
	MaskSize    int
}

type TimeoutsConfig struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}
