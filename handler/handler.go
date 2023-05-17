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
	api.GET(StationRouteWithParams, h.QueryStation)
	api.POST(StationRoute, h.SearchStations)
	api.POST(HistoricalRoute, h.Historical)
	// -- Reset password
	api.POST(ResetPasswordRoute, h.RequestPasswordReset)
	api.PUT(ResetPasswordRoute, h.ResetPassword)
	// Auth Req
	authReq := api.Group(RootRoute, h.Authorize)
	authReq.Any(EchoRoute, h.Echo)
	// -- Account
	authReq.GET(AccountRoute, h.QueryAccount)
	authReq.PATCH(AccountRoute, h.UpdateAccount)
	// -- Alerts
	authReq.PUT(AlertRoute, h.CreateAlert)
	authReq.PATCH(AlertRoute, h.UpdateAlert)
	authReq.DELETE(AlertRouteWithParam, h.DeleteAlert)
	authReq.GET(AlertRouteWithParam, h.QueryAlert)
	authReq.POST(AlertRoute, h.SearchAlerts)
	// -- Stations
	authReq.PUT(StationRoute, h.CreateStation)
	authReq.DELETE(StationRouteWithParams, h.DeleteStation)
	authReq.PATCH(StationRoute, h.UpdateStation)
	// -- Sensors
	authReq.PUT(SensorRoute, h.CreateSensors)
	authReq.DELETE(SensorRoute, h.DeleteSensors)
	// -- Confirm account
	authReq.POST(ConfirmAccountRoute, h.RequestConfirmAccount)
	api.PUT(ConfirmAccountRoute, h.ConfirmAccount)
	return h
}
