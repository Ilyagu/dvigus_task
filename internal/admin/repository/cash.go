package repository

import (
	"sync"
	"time"
)

type RateByIP struct {
	Count int
	Time  time.Time
}

var (
	Mutex = sync.RWMutex{}
	Repo  = make(map[string]RateByIP)
)
