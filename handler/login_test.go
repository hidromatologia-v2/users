package handler

import (
	"net/http"
	"testing"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run("Valid credentials", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.Register(u))
		token := expect.
			POST(LoginRoute).
			WithJSON(u).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			Value("message").
			String().Raw()
		user, aErr := h.Controller.Authorize(token)
		assert.Nil(tt, aErr)
		assert.NotNil(tt, user)
	})
	t.Run("Invalid credentials", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		expect.
			POST(LoginRoute).
			WithJSON(u).
			Expect().
			Status(http.StatusUnauthorized)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			POST(LoginRoute).
			WithBytes([]byte("[")).
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusBadRequest)
	})
}
