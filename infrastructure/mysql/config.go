package mysql

import (
	"context"

	"dario.cat/mergo"
	"github.com/usk81/aveo"
)

type (
	// Config ...
	Config struct {
		DBName               string `env:"DATABASE_NAME"`
		DBHost               string `env:"DATABASE_HOST"`
		DBPort               string `env:"DATABASE_PORT"`
		DBUser               string `env:"DATABASE_USER"`
		DBPassword           string `env:"DATABASE_PASSWORD"`
		DBMaxIdleConnections int    `env:"DATABASE_MAX_IDLE_CONNECTIONS"`
		DBMaxOpenConnections int    `env:"DATABASE_MAX_OPEN_CONNECTIONS"`
	}
)

// NewConfig creates a new Config
func NewConfig(ctx context.Context, env aveo.Env, dt string) Config {
	var rc Config
	_ = aveo.Process(env, ctx, dt+"_", &rc)
	var cf Config
	_ = aveo.Process(env, ctx, "", &cf)
	_ = mergo.Merge(&cf, rc)
	return cf
}
