package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/connection"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/models/common/random"
)

func defaultHandler(t *testing.T) (expect *httpexpect.Expect, h *Handler, stationName string, closeFunc func()) {
	opts := models.Options{
		Database: postgres.NewDefault(),
	}
	c := models.NewController(&opts)
	stationName = random.String()
	h = New(c, connection.NewProducer(t, stationName))
	server := httptest.NewServer(h)
	expect = httpexpect.Default(t, server.URL+APIRoute)
	return expect, h, stationName, func() {
		server.Close()
	}
}
