package config

import (
	"errors"
	"os"
)

type ExporterConfig struct {
	Token   string
	Region  string
}

func Parse() (*ExporterConfig, error) {
	conf := &ExporterConfig{}

	token, ok := os.LookupEnv("SELECTEL_TOKEN")
	if !ok {
		return nil, errors.New("env variable SELECTEL_TOKEN is required")
	}
	conf.Token = token

	region, ok := os.LookupEnv("SELECTEL_REGION")
	if !ok {
		return nil, errors.New("env variable SELECTEL_REGION is required")
	}
	conf.Region = region

	return conf, nil
}
