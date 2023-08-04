package handler

import (
	"errors"
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
	session, sessionExists := ctx.Get(SessionVariable)
	var (
		station *tables.Station
		qErr    error
	)
	if sessionExists {
		station, qErr = h.Controller.QueryStationAPIKey(session.(*tables.User), &s)
		if errors.Is(qErr, models.ErrUnauthorized) {
			station, qErr = h.Controller.QueryStationNoAPIKey(&s)
		}
	} else {
		station, qErr = h.Controller.QueryStationNoAPIKey(&s)
	}
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

func (h *Handler) CreateStation(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var s tables.Station
	bErr := ctx.Bind(&s)
	if bErr != nil {
		return
	}
	station, cErr := h.Controller.CreateStation(session, &s)
	if cErr == nil {
		ctx.JSON(http.StatusCreated, station)
		return
	}
	if errors.Is(cErr, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: cErr.Error()})
}

func (h *Handler) DeleteStation(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var (
		s    tables.Station
		pErr error
	)
	s.UUID, pErr = uuid.FromString(ctx.Param(UUIDParam))
	if pErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	dErr := h.Controller.DeleteStation(session, &s)
	if dErr == nil {
		ctx.JSON(http.StatusOK, SucceedResponse)
		return
	}
	if errors.Is(dErr, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: dErr.Error()})
}

func (h *Handler) UpdateStation(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var s tables.Station
	bErr := ctx.Bind(&s)
	if bErr != nil {
		return
	}
	uErr := h.Controller.UpdateStation(session, &s)
	if uErr == nil {
		ctx.JSON(http.StatusOK, SucceedResponse)
		return
	}
	if errors.Is(uErr, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: uErr.Error()})
}

func (h *Handler) CreateSensors(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var station tables.Station
	bErr := ctx.Bind(&station)
	if bErr != nil {
		return
	}
	err := h.Controller.AddSensors(session, &station, station.Sensors)
	if err == nil {
		ctx.JSON(http.StatusOK, SucceedResponse)
		return
	}
	if errors.Is(err, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: err.Error()})
}

func (h *Handler) DeleteSensors(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var station tables.Station
	bErr := ctx.Bind(&station)
	if bErr != nil {
		return
	}
	err := h.Controller.DeleteSensors(session, &station, station.Sensors)
	if err == nil {
		ctx.JSON(http.StatusOK, SucceedResponse)
		return
	}
	if errors.Is(err, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: err.Error()})
}
