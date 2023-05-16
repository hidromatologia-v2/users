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
