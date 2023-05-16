package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
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

func (h *Handler) SearchStations(ctx *gin.Context) {
	var filter models.Filter[tables.Station]
	bErr := ctx.Bind(&filter)
	if bErr != nil {
		return
	}
	result, sErr := h.Controller.QueryManyStation(&filter)
	if sErr == nil {
		ctx.JSON(http.StatusOK, result)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: sErr.Error()})
}

func (h *Handler) Historical(ctx *gin.Context) {
	var hFilter models.HistoricalFilter
	bErr := ctx.Bind(&hFilter)
	if bErr != nil {
		return
	}
	registries, hErr := h.Controller.Historical(&hFilter)
	if hErr == nil {
		ctx.JSON(http.StatusOK, registries)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: hErr.Error()})
}
