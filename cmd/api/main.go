package main

import (
	"dvigus_task/config/api_config"

	"dvigus_task/internal/pkg/logger"
	middleware "dvigus_task/internal/pkg/middlewares"
	user_delivery "dvigus_task/internal/users/delivery"
	user_usecase "dvigus_task/internal/users/usecase"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	api_config.SetConfig()

	// logger
	logger := logger.NewLogger(api_config.Logger)

	// router
	router := mux.NewRouter()

	// usecase
	uu := user_usecase.NewUserUsecase(logger)

	// middlewars
	loggerMw := middleware.LoggerMiddleware{Logger: logger}
	rateLimitsMw := middleware.RateLimitMiddleware{Repo: make(map[string]middleware.RateByIP)}

	// delivery
	user_delivery.SetUserRouting(router, uu, &loggerMw, &rateLimitsMw, logger)

	srv := &http.Server{
		Handler:      router,
		Addr:         api_config.App.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}
	logger.Infof("starting server at %s\n", srv.Addr)

	logger.Fatal(srv.ListenAndServe())
}
