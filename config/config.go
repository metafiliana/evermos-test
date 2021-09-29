package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	// [db cockroach configuration]
	DBHost                  string `envconfig:"DB_HOST" default:"localhost"`
	DBPort                  string `envconfig:"DB_PORT" default:"3306"`
	DBUserName              string `envconfig:"DB_USERNAME" default:"root"`
	DBName                  string `envconfig:"DB_NAME" default:"evermos-test"`
	DBPass                  string `envconfig:"DB_PASS" default:"root"`
	DBLogMode               bool   `envconfig:"DB_LOG_MODE" default:"true"`
	DBMaxIdleConnection     int    `envconfig:"DB_MAX_IDLE_CONNECTION" default:"5"`
	DBMaxOpenConnection     int    `envconfig:"DB_MAX_OPEN_CONNECTION" default:"10"`
	DBMaxLifetimeConnection int    `envconfig:"DB_MAX_LIFETIME_CONNECTION" default:"10"`
	RedisAddress            string `envconfig:"DB_HOST" default:"localhost:6379"`
	RedisName               string `envconfig:"DB_HOST" default:"evermos-test"`
	RestPort                string `envconfig:"REST_PORT" default:":8080"`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
