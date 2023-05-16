package handler

import (
	"net/http"
	"testing"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/hidromatologia-v2/users/common/headers"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	t.Run("Valid Credentials", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		user, aErr := h.Controller.Authenticate(u)
		assert.Nil(tt, aErr)
		token := h.Controller.JWT.New(user.Claims())
		expect.
			GET(EchoRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK)
	})
	t.Run("Invalid Token", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			GET(EchoRoute).
			WithHeader("Authorization", "Bearer INVALID").
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("Invalid Credentials", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.UUID = uuid.NewV4()
		token := h.Controller.JWT.New(u.Claims())
		expect.
			GET(EchoRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusUnauthorized)
	})
}
