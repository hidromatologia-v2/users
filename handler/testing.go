package handler

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/connection"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
	redis_v9 "github.com/redis/go-redis/v9"
)

func defaultHandler(t *testing.T) (expect *httpexpect.Expect, h *Handler, stationName string, closeFunc func()) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	pCache := cache.Redis(&redis_v9.Options{Addr: "127.0.0.1:6379", DB: 0})
	eCache := cache.Redis(&redis_v9.Options{Addr: "127.0.0.1:6379", DB: 1})
	opts := models.Options{
		Database:      postgres.NewDefault(),
		JWTSecret:     []byte("MY_SECRET"),
		PasswordCache: pCache,
		EmailCache:    eCache,
	}
	c := models.NewController(&opts)
	stationName = random.String()
	h = New(c, connection.NewProducer(t, stationName, random.String()[:64]))
	server := httptest.NewServer(h)
	expect = httpexpect.Default(t, server.URL)
	return expect, h, stationName, func() {
		server.Close()
	}
}
