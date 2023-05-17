package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (h *Handler) RequestPasswordReset(ctx *gin.Context) {
	defer ctx.JSON(http.StatusCreated, SucceedResponse) // Prevent email enumeration
	var u tables.User
	bErr := ctx.Bind(&u)
	if bErr != nil {
		return
	}
	var user tables.User
	fErr := h.Controller.DB.Where("email = ?", u.Email).First(&user).Error
	if fErr != nil {
		return
	}
	resetCode, rErr := h.Controller.RequestResetPassword(&user)
	if rErr != nil {
		return
	}
	msg := &tables.Message{
		Type:      tables.Email,
		Recipient: *user.Email,
		Subject:   "Password reset code",
		Body:      fmt.Sprintf("Your reset code is %s", resetCode),
	}
	var buffer bytes.Buffer
	eErr := json.NewEncoder(&buffer).Encode(msg)
	if eErr != nil {
		return
	}
	pErr := h.MessageProducer.Produce(buffer.Bytes())
	if pErr != nil {
		return
	}
	// TODO: Sure we want to handle errors this way?
}

type ResetRequest struct {
	ResetCode   string `json:"resetCode"`
	NewPassword string `json:"newPassword"`
}

func (h *Handler) ResetPassword(ctx *gin.Context) {
	var req ResetRequest
	bErr := ctx.Bind(&req)
	if bErr != nil {
		return
	}
	rErr := h.Controller.ResetPassword(req.ResetCode, req.NewPassword)
	if rErr == nil {
		ctx.JSON(http.StatusCreated, SucceedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, rErr)
}

func (h *Handler) RequestConfirmAccount(ctx *gin.Context) {
	session := ctx.MustGet(SessionVariable).(*tables.User)
	confirmCode, rErr := h.Controller.RequestConfirmation(session)
	if rErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: rErr.Error()})
		return
	}
	msg := &tables.Message{
		Type:      tables.Email,
		Recipient: *session.Email,
		Subject:   "Confirm account code",
		Body:      fmt.Sprintf("Your confirmation code is %s", confirmCode),
	}
	var buffer bytes.Buffer
	eErr := json.NewEncoder(&buffer).Encode(msg)
	if eErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: eErr.Error()})
		return
	}
	pErr := h.MessageProducer.Produce(buffer.Bytes())
	if pErr == nil {
		ctx.JSON(http.StatusCreated, SucceedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: pErr.Error()})
}

type ConfirmRequest struct {
	ConfirmCode string `json:"confirmCode"`
}

func (h *Handler) ConfirmAccount(ctx *gin.Context) {
	var req ConfirmRequest
	bErr := ctx.Bind(&req)
	if bErr != nil {
		return
	}
	cErr := h.Controller.ConfirmAccount(req.ConfirmCode)
	if cErr == nil {
		ctx.JSON(http.StatusCreated, SucceedResponse)
		return
	}
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, Response{Message: cErr.Error()})
}
