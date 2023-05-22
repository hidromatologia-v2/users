package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/tables"
)

func (h *Handler) Login(ctx *gin.Context) {
	var creds tables.User
	bErr := ctx.Bind(&creds)
	if bErr != nil {
		return
	}
	user, aErr := h.Controller.Authenticate(&creds)
	if aErr == nil {
		token := h.Controller.JWT.New(user.Claims())
		ctx.JSON(http.StatusCreated, Response{Message: token})
		return
	}
	if errors.Is(aErr, models.ErrUnauthorized) {
		ctx.JSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.JSON(http.StatusInternalServerError, Response{Message: aErr.Error()})
}
