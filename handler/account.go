package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
