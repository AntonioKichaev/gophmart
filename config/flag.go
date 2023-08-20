package config

import "flag"

func parseFlag(e *EnvSetting) {
	flag.StringVar(&e.RunAddress, "a", "localhost:8080", "where agent wil send request")
	flag.StringVar(&e.DatabaseURI, "d", "postgres://anton:!anton321@localhost:5435/mart?sslmode=disable", "where agent wil send request")
	flag.StringVar(&e.AccrualSystemAddress, "r", "http://localhost:8088", "where agent wil send request")
	flag.Parse()
}
