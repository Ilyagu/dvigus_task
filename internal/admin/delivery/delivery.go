package delivery

import (
	"dvigus_task/config/api_config"
	"dvigus_task/internal/admin/repository"
	"dvigus_task/internal/pkg/logger"
	middleware "dvigus_task/internal/pkg/middlewares"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"net/http"

	"github.com/c-robinson/iplib"
	"github.com/gorilla/mux"
)

type AdminDelivery struct {
	logger *logger.Logger
}

func SetAdminRouting(
	router *mux.Router,
	loggerMw *middleware.LoggerMiddleware,
	rateLimitsMw *middleware.RateLimitMiddleware,
	logger *logger.Logger,
) {
	adminDelivery := &AdminDelivery{
		logger: logger,
	}

	// public API
	adminPublicAPI := router.PathPrefix("/").Subrouter()

	// установка мидлвар на методы публичные(не стал выносить в другую директорию, так как апи маленькое)
	adminPublicAPI.Use(loggerMw.WithLogging, middleware.CheckRateLimits(logger))

	adminPublicAPI.HandleFunc("/status", adminDelivery.status).Methods(http.MethodGet)

	// private API
	adminPrivateAPI := router.PathPrefix("/admin").Subrouter()

	// установка мидлвар на методы админа
	adminPrivateAPI.Use(loggerMw.WithLogging)

	adminPrivateAPI.HandleFunc("/reset", adminDelivery.reset).Methods(http.MethodGet)
}

func (ud *AdminDelivery) status(w http.ResponseWriter, r *http.Request) {
	// выдача обычного статического контента
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1>Hello</h1>"))
}

func (ud *AdminDelivery) reset(w http.ResponseWriter, r *http.Request) {
	// текущее время
	timeNow := time.Now()

	// вытаскиваем из квери стринга подсеть или же IP
	subNetOrIP := r.URL.Query().Get("subnet")
	if subNetOrIP == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	match, _ := regexp.MatchString(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`, subNetOrIP)
	if !match {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// убираем лишний порт
	subNetOrIP = strings.Split(subNetOrIP, ":")[0]
	// пропускаем IP через заданную маску(дефолт 24) и вычисляем префикс
	n := iplib.NewNet4(net.ParseIP(subNetOrIP), api_config.RateLimit.MaskSize)
	subNet := n.IP().String()

	// блочим запись в мапу по префиксу
	repository.Mutex.RLock()
	if _, ok := repository.Repo[subNet]; !ok {
		// добавляем префикс, если его не было ранее
		repository.Repo[subNet] = repository.RateByIP{
			Time:    timeNow,
			Count:   1,
			TurnOff: true,
		}
	}

	// сбрасываем лимит по префиксу
	rl := repository.Repo[subNet]
	rl.Count = 0
	rl.Time = timeNow
	rl.TurnOff = true
	repository.Repo[subNet] = rl
	repository.Mutex.RUnlock()

	fmt.Println("RESET", rl, subNet)

	w.WriteHeader(http.StatusOK)
}
