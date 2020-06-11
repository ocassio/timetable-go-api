package config

import (
	"os"
	"time"
)

type AppConfig struct {
	Address              string
	TimetableUrl         string
	RequestTimeout       time.Duration
	CacheTimeout         time.Duration
	CacheCleanupInterval time.Duration
}

var Config AppConfig

func init() {
	Config = AppConfig{
		Address:      ":8080",
		TimetableUrl: "https://www.tolgas.ru/services/raspisanie/",

		RequestTimeout:       20,
		CacheTimeout:         24 * 60,
		CacheCleanupInterval: 4 * 60,
	}

	addressEnv, present := os.LookupEnv("ADDRESS")
	if present {
		Config.Address = addressEnv
	}

	timetableUrlEnv, present := os.LookupEnv("TIMETABLE_URL")
	if present {
		Config.TimetableUrl = timetableUrlEnv
	}

	requestTimeoutEnv, present := getDurationVariable("REQUEST_TIMEOUT")
	if present {
		Config.RequestTimeout = *requestTimeoutEnv
	}

	cacheTimeoutEnv, present := getDurationVariable("CACHE_TIMEOUT")
	if present {
		Config.CacheTimeout = *cacheTimeoutEnv
	}

	cacheCleanupIntervalEnv, present := getDurationVariable("CACHE_CLEANUP_INTERVAL")
	if present {
		Config.CacheCleanupInterval = *cacheCleanupIntervalEnv
	}
}

func getDurationVariable(name string) (*time.Duration, bool) {
	variable, present := os.LookupEnv(name)
	if present {
		duration, err := time.ParseDuration(variable)
		if err != nil {
			panic(err)
		}
		return &duration, true
	}

	return nil, false
}
