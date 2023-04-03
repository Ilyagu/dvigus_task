package delivery

import (
	"dvigus_task/internal/pkg/logger"
	middleware "dvigus_task/internal/pkg/middlewares"
	"dvigus_task/internal/users/usecase"
	"path"
	"path/filepath"
	"runtime"

	"net/http"

	"github.com/gorilla/mux"
)

type UserDelivery struct {
	logger *logger.Logger
}

func SetUserRouting(
	router *mux.Router,
	uu usecase.IUserUsecase,
	loggerMw *middleware.LoggerMiddleware,
	rateLimitsMw *middleware.RateLimitMiddleware,
	logger *logger.Logger,
) {
	userDelivery := &UserDelivery{
		logger: logger,
	}

	// public API
	userPublicAPI := router.PathPrefix("/").Subrouter()

	userPublicAPI.Use(loggerMw.WithLogging, rateLimitsMw.CheckRateLimits)

	userPublicAPI.HandleFunc("/status", userDelivery.index).Methods(http.MethodGet)
}

func (ud *UserDelivery) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte("<h1>Hello</h1>"))
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
