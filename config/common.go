package config

import (
	"time"
)

type AppConfig struct {
	Port                  string
	WellSimilarityPercent int
	ExpirationCookieTime  time.Duration
}

type RateLimitConfig struct {
	Period      time.Duration
	MaxRequests int
	Cooldown    time.Duration
	Status      bool
	MaskSize    int
}

type TimeoutsConfig struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}
