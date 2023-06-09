package handler

import (
	"net/http"
	"testing"

	"github.com/hidromatologia-v2/models"
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

func TestUpdateAlert(t *testing.T) {
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
		a := tables.RandomAlert(u, &sensor)
		assert.Nil(tt, h.Controller.DB.Create(a).Error)
		newName := random.String()
		expect.
			PATCH(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Alert{
				Model:    a.Model,
				UserUUID: u.UUID,
				Name:     &newName,
				Enabled:  &tables.True,
			}).
			Expect().
			Status(http.StatusOK)
		var alert tables.Alert
		assert.Nil(tt, h.Controller.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
		assert.True(tt, *alert.Enabled)
		assert.NotEqual(tt, *a.Name, *alert.Name)
	})
	t.Run("Other user alert", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		sensor := s.Sensors[0]
		a := tables.RandomAlert(u2, &sensor)
		assert.Nil(tt, h.Controller.DB.Create(a).Error)
		newName := random.String()
		expect.
			PATCH(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Alert{
				Model:    a.Model,
				UserUUID: u2.UUID,
				Name:     &newName,
				Enabled:  &tables.True,
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PATCH(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestDeleteAlert(t *testing.T) {
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
		a := tables.RandomAlert(u, &sensor)
		assert.Nil(tt, h.Controller.DB.Create(a).Error)
		expect.
			DELETE(AlertRoute+"/"+a.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK)
		var alert tables.Alert
		assert.NotNil(tt, h.Controller.DB.Where("uuid = ?", a.UUID).First(&alert).Error)
	})
	t.Run("Other user alert", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		sensor := s.Sensors[0]
		a := tables.RandomAlert(u2, &sensor)
		assert.Nil(tt, h.Controller.DB.Create(a).Error)
		expect.
			DELETE(AlertRoute+"/"+a.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusUnauthorized)
	})
	t.Run("Invalid UUID", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			DELETE(AlertRoute+"/INVALID").
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestQueryAlert(t *testing.T) {
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
		a := tables.RandomAlert(u, &sensor)
		assert.Nil(tt, h.Controller.DB.Create(a).Error)
		var alert tables.Alert
		expect.
			GET(AlertRoute+"/"+a.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().Decode(&alert)
		assert.Equal(tt, a.UUID, alert.UUID)

	})
	t.Run("Other user alert", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		sensor := s.Sensors[0]
		a := tables.RandomAlert(u2, &sensor)
		assert.Nil(tt, h.Controller.DB.Create(a).Error)
		expect.
			GET(AlertRoute+"/"+a.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusUnauthorized)
	})
	t.Run("Invalid UUID", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			GET(AlertRoute+"/INVALID").
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestSearchAlerts(t *testing.T) {
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
		// -- Create alerts
		for i := 0; i < 10; i++ {
			a := tables.RandomAlert(u, &sensor)
			assert.Nil(tt, h.Controller.DB.Create(a).Error)
		}
		// --
		var results models.Results[tables.Alert]
		expect.
			POST(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(models.Filter[tables.Alert]{
				Page:     1,
				PageSize: 100,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(&results)
		assert.Equal(tt, 10, results.Count)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			POST(AlertRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
}
