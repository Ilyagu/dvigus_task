package delivery

import (
	"dvigus_task/config/api_config"
	"dvigus_task/internal/admin/repository"
	"dvigus_task/internal/pkg/logger"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

type MyHandler http.Handler

func TestReset(t *testing.T) {
	// t.Parallel()

	logger := logger.NewLogger(api_config.Logger)
	api_config.RateLimit.MaskSize = 24
	// api_config.RateLimit.Cooldown = time.Second * 2
	// api_config.RateLimit.Period = time.Second * 3
	// api_config.RateLimit.MaxRequests = 1

	adminHandler := AdminDelivery{logger}
	t.Run("200", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/admin/reset?subnet=127.0.0.1", nil)
		w := httptest.NewRecorder()

		adminHandler.reset(w, r)

		statOfSubnet := repository.Repo["127.0.0.0"]
		fmt.Println(statOfSubnet)
		assert.Equal(t, w.Code, 200)
	})
	t.Run("400", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/admin/reset?subnet=127.0.0.deaf", nil)
		w := httptest.NewRecorder()

		adminHandler.reset(w, r)

		statOfSubnet := repository.Repo["127.0.0.0"]
		fmt.Println(statOfSubnet)
		assert.Equal(t, w.Code, 400)
	})
}

func TestStatus(t *testing.T) {
	logger := logger.NewLogger(api_config.Logger)

	adminHandler := AdminDelivery{logger}
	r := httptest.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()

	adminHandler.status(w, r)

	assert.Equal(t, string(w.Body.Bytes()), "<h1>Hello</h1>")
	assert.Equal(t, w.Code, 200)
}
