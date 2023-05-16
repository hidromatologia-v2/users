package handler

import (
	"net/http"
	"testing"

	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/hidromatologia-v2/users/common/headers"
	"github.com/stretchr/testify/assert"
)

func TestCreateAlert(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		sensor := s.Sensors[0]
		aName := random.String()
		condition := tables.Lt
		value := random.Float(1000.0)
		expect.
			PUT(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Alert{
				UserUUID:   u.UUID,
				SensorUUID: sensor.UUID,
				Name:       &aName,
				Condition:  &condition,
				Value:      &value,
				Enabled:    &tables.True,
			}).
			Expect().
			Status(http.StatusCreated)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PUT(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
}
