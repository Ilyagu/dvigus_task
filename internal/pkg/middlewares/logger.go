package middleware

import (
	"dvigus_task/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	Logger *logger.Logger
}

func (lm LoggerMiddleware) WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := lm.Logger.WithFields(logrus.Fields{
			"url":    r.URL,
			"method": r.Method,
		})

		log.Infof("Start request")
		start := time.Now()

		h.ServeHTTP(w, r)

		log.WithFields(logrus.Fields{
			"duration": time.Since(start),
		}).Info("End request")
	})
}
