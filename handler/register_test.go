package handler

import (
	"net/http"
	"testing"

	"github.com/hidromatologia-v2/models/tables"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		expect.
			PUT(RegisterRoute).
			WithJSON(u).
			Expect().
			Status(http.StatusCreated)
		var user tables.User
		assert.Nil(tt, h.Controller.DB.Where("username = ?", u.Username).First(&user).Error)
		assert.Equal(tt, u.Username, user.Username)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			PUT(RegisterRoute).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("{")).
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("Repeated User", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		expect.
			PUT(RegisterRoute).
			WithJSON(u).
			Expect().
			Status(http.StatusCreated)
		var user tables.User
		assert.Nil(tt, h.Controller.DB.Where("username = ?", u.Username).First(&user).Error)
		assert.Equal(tt, u.Username, user.Username)
		// User 2
		expect.
			PUT(RegisterRoute).
			WithJSON(u).
			Expect().
			Status(http.StatusInternalServerError)

	})
}
