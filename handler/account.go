package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/tables"
)

func (h *Handler) QueryAccount(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	user, qErr := h.Controller.QueryAccount(session)
	if qErr == nil {
		ctx.JSON(http.StatusOK, user)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: qErr.Error()})
}

type UpdateAccountRequest struct {
	User        tables.User `json:"user"`
	OldPassword string      `json:"oldPassword"`
}

func (h *Handler) UpdateAccount(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	var req UpdateAccountRequest
	bErr := ctx.Bind(&req)
	if bErr != nil {
		return
	}
	uErr := h.Controller.UpdateAccount(session, &req.User, req.OldPassword)
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
