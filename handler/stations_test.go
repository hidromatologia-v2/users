package handler

import (
	"net/http"
	"testing"
	"time"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/random"
	"github.com/hidromatologia-v2/models/tables"
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
			GET(QueryStationRoute + "/" + s.UUID.String()).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Decode(&station)
		assert.Equal(tt, s.UUID, station.UUID)
	})
	t.Run("Invalid UUID", func(tt *testing.T) {
		expect, h, _, closeFunc := defaultHandler(tt)
		defer h.Close()
		defer closeFunc()
		expect.
			GET(QueryStationRoute + "/INVALID").
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
			POST(SearchStationsRoute).
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
			POST(SearchStationsRoute).
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
