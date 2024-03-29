package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
	"github.com/hidromatologia-v2/users/common/headers"
	"github.com/stretchr/testify/assert"
)

func TestQueryStation(t *testing.T) {
	t.Run("Existing Station", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		var station tables.Station
		expect.
			GET(StationRoute + "/" + s.UUID.String()).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(&station)
		assert.Equal(tt, s.UUID, station.UUID)
	})

	t.Run("With Session", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		var station tables.Station
		expect.
			GET(StationRoute+"/"+s.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(&station)
		assert.Equal(tt, s.UUID, station.UUID)
		assert.NotNil(tt, station.APIKeyJSON)
		assert.Equal(tt, s.APIKey, *station.APIKeyJSON)
	})
	t.Run("With Session But Not owner", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u2.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		var station tables.Station
		expect.
			GET(StationRoute+"/"+s.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(&station)
		assert.Equal(tt, s.UUID, station.UUID)
		assert.Nil(tt, station.APIKeyJSON)
	})
	t.Run("Invalid UUID", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			GET(StationRoute + "/INVALID").
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestSearchStations(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			POST(StationRoute).
			WithJSON(models.Filter[tables.Station]{
				Page:     1,
				PageSize: 100,
				Target: tables.Station{
					Name: s.Name,
				},
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Value("count").
			Number().IsEqual(1)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			POST(StationRoute).
			WithBytes([]byte("[")).
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestHistorical(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		sensor := s.Sensors[0]
		from := time.Now()
		for i := 0; i < 10; i++ {
			sr := tables.SensorRegistry{
				SensorUUID: sensor.UUID,
				Value:      random.Float(1000.0),
			}
			assert.Nil(tt, h.Controller.DB.Create(&sr).Error)
		}
		to := time.Now()
		expect.
			POST(HistoricalRoute).
			WithJSON(models.HistoricalFilter{
				SensorUUID: sensor.UUID,
				From:       &from,
				To:         &to,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Array().
			Length().IsEqual(10)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			POST(HistoricalRoute).
			WithBytes([]byte("[")).
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestCreateStation(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		var station tables.Station
		expect.
			PUT(StationRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(s).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			Decode(&station)
	})
	t.Run("Not confirmed", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		expect.
			PUT(StationRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(s).
			Expect().
			Status(http.StatusUnauthorized)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PUT(StationRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestDeleteStation(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			DELETE(StationRoute+"/"+s.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusOK)
	})
	t.Run("Other user station", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u2)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			DELETE(StationRoute+"/"+s.UUID.String()).
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusUnauthorized)
	})
	t.Run("Invalid UUID", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			DELETE(StationRoute+"/INVALID").
			WithHeader("Authorization", headers.Authorization(token)).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func TestUpdateStation(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		newName := random.String()
		expect.
			PATCH(StationRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Station{
				Model: s.Model,
				Name:  &newName,
			}).
			Expect().
			Status(http.StatusOK)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PATCH(StationRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("Other user station", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u2)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		newName := random.String()
		expect.
			PATCH(StationRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Station{
				Model: s.Model,
				Name:  &newName,
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func TestCreateSensors(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		s.Sensors = nil
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			PUT(SensorRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Station{
				Model: s.Model,
				Sensors: []tables.Sensor{
					{
						Type: random.String(),
					},
					{
						Type: random.String(),
					},
					{
						Type: random.String(),
					},
					{
						Type: random.String(),
					},
				},
			}).
			Expect().
			Status(http.StatusOK)
		var station tables.Station
		assert.Nil(tt, h.Controller.DB.Where("uuid = ?", s.UUID).Preload("Sensors").First(&station).Error)
		assert.Len(tt, station.Sensors, 4)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PUT(SensorRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("Other user station", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u2)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			PUT(SensorRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Station{
				Model: s.Model,
				Sensors: []tables.Sensor{
					{
						Type: random.String(),
					},
					{
						Type: random.String(),
					},
					{
						Type: random.String(),
					},
					{
						Type: random.String(),
					},
				},
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func TestDeleteSensors(t *testing.T) {
	t.Run("Valid", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			PATCH(SensorRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Station{
				Model:   s.Model,
				Sensors: s.Sensors,
			}).
			Expect().
			Status(http.StatusOK)
		var station tables.Station
		assert.Nil(tt, h.Controller.DB.Where("uuid = ?", s.UUID).Preload("Sensors").First(&station).Error)
		assert.Len(tt, station.Sensors, 0)
	})
	t.Run("Invalid JSON", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		u.Confirmed = &tables.True
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		token := h.Controller.JWT.New(u.Claims())
		expect.
			PATCH(SensorRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithHeader("Content-Type", "application/json").
			WithBytes([]byte("[")).
			Expect().
			Status(http.StatusBadRequest)
	})
	t.Run("Other user station", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		u := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u).Error)
		u2 := tables.RandomUser()
		assert.Nil(tt, h.Controller.DB.Create(u2).Error)
		token := h.Controller.JWT.New(u.Claims())
		s := tables.RandomStation(u2)
		assert.Nil(tt, h.Controller.DB.Create(s).Error)
		expect.
			PATCH(SensorRoute).
			WithHeader("Authorization", headers.Authorization(token)).
			WithJSON(tables.Station{
				Model:   s.Model,
				Sensors: s.Sensors,
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})
}
