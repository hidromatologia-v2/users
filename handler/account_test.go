package handler

import (
	"net/http"
	"testing"

	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/hidromatologia-v2/users/common/headers"
	"github.com/stretchr/testify/assert"
)

func TestQueryAccount(t *testing.T) {
	t.Run("Account Details", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		var user tables.User
		expect.
			GET(AccountRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(&user)
		assert.Equal(tt, u.UUID, user.UUID)
	})
}

func TestUpdateAccount(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		newName := random.String()
		expect.
			PATCH(AccountRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(UpdateAccountRequest{
				User: tables.User{
					Name: &newName,
				},
				OldPassword: u.Password,
			}).
			Expect().
			Status(http.StatusOK)
		var user tables.User
		assert.Nil(tt, h.Controller.DB.Where("uuid = ?", u.UUID).First(&user).Error)
		assert.NotEqual(tt, *u.Name, *user.Name)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PATCH(AccountRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("Invalid old password", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		newName := random.String()
		expect.
			PATCH(AccountRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(UpdateAccountRequest{
				User: tables.User{
					Name: &newName,
				},
				OldPassword: random.String(),
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func TestRequestResetPassword(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		expect.
			POST(ResetPasswordRoute).
			WithJSON(u).
			Expect().
			Status(http.StatusCreated)
	})
}

func TestResetPassword(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		resetCode, rErr := h.Controller.RequestResetPassword(u)
		assert.Nil(tt, rErr)
		newPassword := random.String()
		expect.
			PUT(ResetPasswordRoute).
			WithJSON(ResetRequest{
				ResetCode:   resetCode,
				NewPassword: newPassword[:72],
			}).
			Expect().
			Status(http.StatusCreated)
	})
	t.Run("Invalid code", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		expect.
			PUT(ResetPasswordRoute).
			WithJSON(ResetRequest{
				ResetCode:   random.String(),
				NewPassword: random.String(),
			}).
			Expect().
			Status(http.StatusInternalServerError)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		expect.
			PUT(ResetPasswordRoute).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestRequestConfirmAccount(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			POST(ConfirmAccountRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusCreated)
	})
}

func TestConfirmAccount(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		confirmCode, cErr := h.Controller.RequestConfirmation(u)
		assert.Nil(tt, cErr)
		expect.
			PUT(ConfirmAccountRoute).
			WithJSON(ConfirmRequest{
				ConfirmCode: confirmCode,
			}).
			Expect().
			Status(http.StatusCreated)
	})
	t.Run("Invalid code", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		expect.
			PUT(ConfirmAccountRoute).
			WithJSON(ConfirmRequest{
				ConfirmCode: random.String(),
			}).
			Expect().
			Status(http.StatusInternalServerError)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		expect.
			PUT(ConfirmAccountRoute).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
}
