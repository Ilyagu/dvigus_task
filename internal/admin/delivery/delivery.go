package delivery

import (
	"dvigus_task/config/api_config"
	"dvigus_task/internal/admin/repository"
	"dvigus_task/internal/pkg/logger"
	middleware "dvigus_task/internal/pkg/middlewares"
	"net"
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

	adminPublicAPI.Use(loggerMw.WithLogging, rateLimitsMw.CheckRateLimits)

	adminPublicAPI.HandleFunc("/status", adminDelivery.status).Methods(http.MethodGet)
	adminPublicAPI.HandleFunc("/reset", adminDelivery.reset).Methods(http.MethodGet)
}

func (ud *AdminDelivery) status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1>Hello</h1>"))
}

func (ud *AdminDelivery) reset(w http.ResponseWriter, r *http.Request) {
	timeNow := time.Now()

	subNetOrIP := r.URL.Query().Get("subnet")

	n := iplib.NewNet4(net.ParseIP(subNetOrIP), api_config.RateLimit.MaskSize)
	subNet := n.IP()

	repository.Mutex.RLock()
	if _, ok := repository.Repo[subNet.String()]; !ok {
		repository.Repo[subNet.String()] = repository.RateByIP{
			Time:  timeNow,
			Count: 1,
		}
	}

	rl := repository.Repo[subNet.String()]
	rl.Count = 0
	rl.Time = timeNow
	repository.Repo[subNet.String()] = rl
	repository.Mutex.RUnlock()
}
