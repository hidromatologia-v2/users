package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models/tables"
)

func (h *Handler) Register(ctx *gin.Context) {
	var user tables.User
	bErr := ctx.Bind(&user)
	if bErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, Response{Message: bErr.Error()})
		return
	}
	rErr := h.Controller.Register(&user)
	if rErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: rErr.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, SucceedResponse)
	ctx.Done()
}
