package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/tables"
	uuid "github.com/satori/go.uuid"
)

func (h *Handler) CreateAlert(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var alert tables.Alert
	bErr := ctx.Bind(&alert)
	if bErr != nil {
		return
	}
	cErr := h.Controller.CreateAlert(session, &alert)
	if cErr == nil {
		ctx.JSON(http.StatusCreated, SucceedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: cErr.Error()})
}

func (h *Handler) UpdateAlert(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var alert tables.Alert
	bErr := ctx.Bind(&alert)
	if bErr != nil {
		return
	}
	cErr := h.Controller.UpdateAlert(session, &alert)
	if cErr == nil {
		ctx.JSON(http.StatusOK, SucceedResponse)
		return
	}
	if errors.Is(cErr, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: cErr.Error()})
}

func (h *Handler) DeleteAlert(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var (
		alert tables.Alert
		pErr  error
	)
	alert.UUID, pErr = uuid.FromString(ctx.Param(UUIDParam))
	if pErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	dErr := h.Controller.DeleteAlert(session, &alert)
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
