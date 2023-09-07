package configs

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

var RunAddress string
var DatabaseURI string
var AccrualSystemAddress string

var (
	runAddressFlag           = flag.String("a", "localhost:8080", "Адрес и порт запуска сервиса")
	databaseURIFlag          = flag.String("d", "host=localhost user=yandex password=yandex dbname=go-diploma-1 sslmode=disable", "Адрес подключения к базе данных")
	accrualSystemAddressFlag = flag.String("r", "", "Адрес системы расчёта начислений")
)

type OsEnvs struct {
	RunAddress           *string `env:"RUN_ADDRESS"`
	DatabaseURI          *string `env:"DATABASE_URI"`
	AccrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Configure() error {
	flag.Parse()

	osEnvs := OsEnvs{}

	err := env.Parse(&osEnvs)
	if err != nil {
		return err
	}

	if osEnvs.RunAddress != nil {
		RunAddress = *osEnvs.RunAddress
	} else {
		RunAddress = *runAddressFlag
	}

	if osEnvs.DatabaseURI != nil {
		DatabaseURI = *osEnvs.DatabaseURI
	} else {
		DatabaseURI = *databaseURIFlag
	}

	if osEnvs.AccrualSystemAddress != nil {
		AccrualSystemAddress = *osEnvs.AccrualSystemAddress
	} else {
		AccrualSystemAddress = *accrualSystemAddressFlag
	}

	return nil
}
