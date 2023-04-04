package main

import (
	"dvigus_task/config/api_config"

	admin_delivery "dvigus_task/internal/admin/delivery"
	"dvigus_task/internal/pkg/logger"
	middleware "dvigus_task/internal/pkg/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	api_config.SetConfig()

	// logger
	logger := logger.NewLogger(api_config.Logger)

	// router
	router := mux.NewRouter()

	// middlewars
	loggerMw := middleware.LoggerMiddleware{Logger: logger}
	rateLimitsMw := middleware.RateLimitMiddleware{}

	// delivery
	admin_delivery.SetAdminRouting(router, &loggerMw, &rateLimitsMw, logger)

	srv := &http.Server{
		Handler:      router,
		Addr:         api_config.App.Port,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}
	logger.Infof("starting server at %s\n", srv.Addr)

	logger.Fatal(srv.ListenAndServe())
}
