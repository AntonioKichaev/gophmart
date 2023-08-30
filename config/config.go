package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/sirupsen/logrus"
)

type EnvSetting struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKeyAuth        string `env:"SecretKeyAuth"`
}

type Config struct {
	env *EnvSetting
}

func New() *Config {
	cfg := &Config{
		env: &EnvSetting{
			RunAddress:           "localhost:8080",
			DatabaseURI:          "postgres://anton:!anton321@localhost:5444/mart?sslmode=disable",
			AccrualSystemAddress: "",
			SecretKeyAuth:        "322",
		},
	}

	parseFlag(cfg.env)
	if err := env.Parse(cfg.env); err != nil {
		logrus.Errorf("parse env err:= %v", err)
	}
	return cfg
}

func (c *Config) GetDatabaseURI() string {
	return c.env.DatabaseURI
}

func (c *Config) GetRunAddress() string {
	return c.env.RunAddress
}

func (c *Config) GetAccrualSystemAddress() string {
	return c.env.AccrualSystemAddress
}

func (c *Config) GetSecretKeyAuth() string {
	return c.env.SecretKeyAuth
}
