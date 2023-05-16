package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/hidromatologia-v2/models"
)

var bearerRegexp = regexp.MustCompile(`(?mi)bearer\s+`)

func (h *Handler) Authorize(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")
	b64Token := bearerRegexp.ReplaceAllString(auth, "")
	token, dErr := base64.StdEncoding.DecodeString(b64Token)
	if dErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user, aErr := h.Controller.Authorize(string(token))
	if aErr == nil {
		ctx.Set(SessionVariable, user)
		ctx.Next()
		return
	}
	if errors.Is(aErr, models.ErrUnauthorized) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: aErr.Error()})
}

func (h *Handler) Echo(ctx *gin.Context) {}
