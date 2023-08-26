package main

import (
	"context"
	"log"
	"os"

	"github.com/hidromatologia-v2/models"
	"github.com/hidromatologia-v2/models/common/cache"
	"github.com/hidromatologia-v2/models/common/postgres"
	"github.com/hidromatologia-v2/users/handler"
	"github.com/memphisdev/memphis.go"
	redis_v9 "github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"github.com/sethvargo/go-envconfig"
)

func logFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func newProducer(config *Config) *memphis.Producer {
	var connOpts []memphis.Option
	if config.Producer.Password != nil {
		connOpts = append(connOpts, memphis.Password(*config.Producer.Password))
	}
	if config.Producer.ConnectionToken != nil {
		connOpts = append(connOpts, memphis.ConnectionToken(*config.Producer.ConnectionToken))
	}
	conn, connErr := memphis.Connect(
		config.Producer.Host,
		config.Producer.Username,
		connOpts...,
	)
	logFatalErr(connErr)
	producer, producerErr := conn.CreateProducer(
		config.Producer.Station,
		config.Producer.Producer+uuid.NewV4().String(),
	)
	logFatalErr(producerErr)
	return producer
}

func main() {
	var config Config
	eErr := envconfig.Process(context.Background(), &config)
	logFatalErr(eErr)
	controllerOpts := models.Options{
		Database: postgres.New(config.Postgres.DSN),
		EmailCache: cache.Redis(&redis_v9.Options{
			Addr: config.EmailCache.Addr,
			DB:   config.EmailCache.DB,
		}),
		PasswordCache: cache.Redis(&redis_v9.Options{
			Addr: config.PasswordCache.Addr,
			DB:   config.PasswordCache.DB,
		}),
		JWTSecret: []byte(config.JWT.Secret),
	}
	producer := newProducer(&config)
	h := handler.New(models.NewController(&controllerOpts), producer)
	rErr := h.Run(os.Args[1:]...)
	if rErr != nil {
		log.Fatal(rErr)
	}
}
