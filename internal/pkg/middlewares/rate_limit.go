package middleware

import (
	"dvigus_task/config/api_config"
	"dvigus_task/internal/admin/repository"
	"dvigus_task/internal/pkg/logger"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/c-robinson/iplib"
)

type RateLimitMiddleware struct {
	Logger *logger.Logger
}

func CheckRateLimits(l *logger.Logger) (mw func(http.Handler) http.Handler) {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// текущее время
			timeNow := time.Now()
			// достаем IP адрес из хедера или из запроса
			ipAddress := r.RemoteAddr
			fwdAddress := r.Header.Get("X-Forwarded-For") // достаем IP
			if fwdAddress != "" {
				// Достаем первый элемент из хедера
				ips := strings.Split(fwdAddress, ", ")
				if len(ips) > 1 {
					ipAddress = ips[0] // Если один элемент
				} else {
					ipAddress = fwdAddress // Если один элемент
				}
			}

			// убираем лишний хост
			ipAddress = strings.Split(ipAddress, ":")[0]
			match, _ := regexp.MatchString(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`, ipAddress)
			if !match {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			n := iplib.NewNet4(net.ParseIP(ipAddress), api_config.RateLimit.MaskSize)
			subNet := n.IP().String() // вытаскиваем подсеть

			// если не было ранее префикса, то добавляем
			if _, ok := repository.Repo[subNet]; !ok {
				repository.Repo[subNet] = repository.RateByIP{
					Time:  timeNow,
					Count: 1,
				}
				h.ServeHTTP(w, r)
				return
			}

			// достаем статистику префикса
			rl := repository.Repo[subNet]
			fmt.Println("RATE LIMIT", rl, subNet, timeNow)
			/// если админ отключил лимит по префиксу просто продолжаем запрос
			if rl.TurnOff {
				h.ServeHTTP(w, r)
				return
			}
			if rl.Count == api_config.RateLimit.MaxRequests {
				// если количество запросов максимум, то выдаем 429 либо сбрасываем после кулдауна
				if timeNow.Sub(rl.Time) > api_config.RateLimit.Cooldown {
					rl.Time = timeNow
					rl.Count = 1
				} else {
					w.WriteHeader(http.StatusTooManyRequests)
					return
				}
			} else {
				// если период прошел, то обнуляем счетчик либо инкрементируем
				// если количество меньше на единицу выставляем текущее время как последнее
				if timeNow.Sub(rl.Time) > api_config.RateLimit.Period {
					rl.Time = timeNow
					rl.Count = 1
				} else {
					if rl.Count == api_config.RateLimit.MaxRequests-1 {
						rl.Time = timeNow
					}
					rl.Count++
				}
			}

			// запись по префиксу
			repository.Repo[subNet] = rl

			h.ServeHTTP(w, r)
		})
	}
}
