package middleware

import (
	"dvigus_task/config/api_config"
	"dvigus_task/internal/admin/repository"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/c-robinson/iplib"
)

type RateLimitMiddleware struct{}

func (rm RateLimitMiddleware) CheckRateLimits(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()
		ipAddress := r.RemoteAddr
		fwdAddress := r.Header.Get("X-Forwarded-For") // достаем IP
		if fwdAddress != "" {
			// Достаем первый элемент из хедера
			ips := strings.Split(fwdAddress, ", ")
			if len(ips) > 1 {
				ipAddress = ips[0]
			} else {
				ipAddress = fwdAddress // Если один элемент
			}
		}

		n := iplib.NewNet4(net.ParseIP(ipAddress), api_config.RateLimit.MaskSize)
		subNet := n.IP()

		if _, ok := repository.Repo[subNet.String()]; !ok {
			repository.Repo[subNet.String()] = repository.RateByIP{
				Time:  timeNow,
				Count: 1,
			}
			h.ServeHTTP(w, r)
			return
		}

		rl := repository.Repo[subNet.String()]
		fmt.Println(rl, timeNow)
		if rl.Count == api_config.RateLimit.MaxRequests {
			if timeNow.Sub(rl.Time) > api_config.RateLimit.Cooldown {
				rl.Time = timeNow
				rl.Count = 1
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		} else {
			if timeNow.Sub(rl.Time) > api_config.RateLimit.Period {
				rl.Time = timeNow
				rl.Count = 1
			} else {
				if rl.Count == api_config.RateLimit.MaxRequests-1 {
					rl.Count++
					rl.Time = timeNow
				} else {
					rl.Count++
				}
			}
		}
		repository.Repo[subNet.String()] = rl

		h.ServeHTTP(w, r)
	})
}
