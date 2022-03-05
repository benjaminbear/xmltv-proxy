package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Days          int
	TimeZone      *time.Location
	DailyDownload bool
	Port          int
	InsecureTLS   bool
}

const FolderEPGData = "epgdata_files"

func ParseEnv() (conf *Config, err error) {
	conf = &Config{
		Days: 7,
	}

	dayStr := os.Getenv("XMLTV_DAYS")
	if dayStr != "" {
		conf.Days, err = strconv.Atoi(dayStr)
		if err != nil {
			return conf, err
		}
	}

	tz := os.Getenv("XMLTV_TIMEZONE")
	if tz == "" {
		conf.TimeZone = time.Now().Local().Location()
	} else {
		conf.TimeZone, err = time.LoadLocation(tz)
		if err != nil {
			return conf, err
		}
	}

	dd := os.Getenv("XMLTV_DAILY_DOWNLOAD")
	if dd != "" {
		conf.DailyDownload, err = strconv.ParseBool(dd)
		if err != nil {
			return conf, err
		}
	}

	port := os.Getenv("XMLTV_PORT")
	if port == "" {
		conf.Port = 8080
	} else {
		conf.Port, err = strconv.Atoi(port)
		if err != nil {
			return conf, err
		}
	}

	ica := os.Getenv("XMLTV_INSECURE_TLS")
	if ica != "" {
		conf.InsecureTLS, err = strconv.ParseBool(ica)
		if err != nil {
			return conf, err
		}
	}

	return conf, err
}
