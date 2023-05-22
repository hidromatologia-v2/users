package main

import "github.com/hidromatologia-v2/models/common/config"

type (
	Config struct {
		EmailCache      config.Redis              `env:",prefix=REDIS_EMAIL_"`
		PasswordCache   config.Redis              `env:",prefix=REDIS_PASSWORD_"`
		config.JWT      `env:",prefix=JWT_"`      // JWT
		config.Producer `env:",prefix=MEMPHIS_"`  // Memphis
		config.Postgres `env:",prefix=POSTGRES_"` // POSTGRESQL
	}
)
