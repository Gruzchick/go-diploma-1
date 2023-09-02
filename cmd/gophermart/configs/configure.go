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
	runAddress           *string `env:"RUN_ADDRESS"`
	databaseURI          *string `env:"DATABASE_URI"`
	accrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Configure() error {
	flag.Parse()

	osEnvs := OsEnvs{}

	err := env.Parse(&osEnvs)
	if err != nil {
		return err
	}

	if osEnvs.runAddress != nil {
		RunAddress = *osEnvs.runAddress
	} else {
		RunAddress = *runAddressFlag
	}

	if osEnvs.databaseURI != nil {
		DatabaseURI = *osEnvs.databaseURI
	} else {
		DatabaseURI = *databaseURIFlag
	}

	if osEnvs.accrualSystemAddress != nil {
		AccrualSystemAddress = *osEnvs.accrualSystemAddress
	} else {
		AccrualSystemAddress = *accrualSystemAddressFlag
	}

	return nil
}
