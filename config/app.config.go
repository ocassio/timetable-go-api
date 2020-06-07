package config

import (
	"os"
	"time"
)

type AppConfig struct {
	Address      string
	TimetableUrl string
	RequestTimeout time.Duration
}

var Config AppConfig

func init() {
	Config = AppConfig {
		Address:      ":8080",
		TimetableUrl: "https://www.tolgas.ru/services/raspisanie/",
		RequestTimeout: 20,
	}

	addressEnv, present := os.LookupEnv("ADDRESS"); if present {
		Config.Address = addressEnv
	}

	timetableUrlEnv, present := os.LookupEnv("TIMETABLE_URL"); if present {
		Config.TimetableUrl = timetableUrlEnv
	}

	requestTimeoutEnv, present := os.LookupEnv("REQUEST_TIMEOUT"); if present {
		timeout, err := time.ParseDuration(requestTimeoutEnv); if err != nil {
			panic(err)
		}
		Config.RequestTimeout = timeout
	}
}
