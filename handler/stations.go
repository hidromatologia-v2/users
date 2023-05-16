package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models/tables"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) QueryStation(ctx *gin.Context) {
	var (
		s    tables.Station
		pErr error
	)
	s.UUID, pErr = uuid.FromString(ctx.Param(UUIDParam))
	if pErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	station, qErr := h.Controller.QueryStationNoAPIKey(&s)
	if qErr == nil {
		ctx.JSON(http.StatusOK, station)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: qErr.Error()})
}
