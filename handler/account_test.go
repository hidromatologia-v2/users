package handler

import (
	"net/http"
	"testing"

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
