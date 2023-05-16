package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
	"github.com/memphisdev/memphis.go"
)

type Handler struct {
	Controller      *models.Controller
	MessageProducer *memphis.Producer
	*gin.Engine
}

func (h *Handler) Close() {
	h.MessageProducer.Destroy()
	h.Controller.Close()
}

func New(c *models.Controller, msgProducer *memphis.Producer) *Handler {
	h := &Handler{
		Controller:      c,
		MessageProducer: msgProducer,
		Engine:          gin.Default(),
	}
	api := h.Group(APIRoute)
	// -- Public
	api.PUT(RegisterRoute, h.Register)
	api.POST(LoginRoute, h.Login)
	api.GET(QueryStationRouteWithParams, h.QueryStation)
	api.POST(SearchStationsRoute, h.SearchStations)
	api.POST(HistoricalRoute, h.Historical)
	// Auth Req
	authReq := api.Group(RootRoute, h.Authorize)
	authReq.Any(EchoRoute, h.Echo)
	// -- Account
	authReq.GET(AccountRoute, h.QueryAccount)
	authReq.PATCH(AccountRoute, h.UpdateAccount)
	// -- Alerts
	authReq.PUT(AlertRoute, h.CreateAlert)
	// -- Stations
	return h
}
