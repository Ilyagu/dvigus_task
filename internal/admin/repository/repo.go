package repository

import (
	"sync"
	"time"
)

type RateByIP struct {
	Count   int
	Time    time.Time
	TurnOff bool
}

var (
	Mutex = sync.RWMutex{}
	Repo  = make(map[string]RateByIP)
)
