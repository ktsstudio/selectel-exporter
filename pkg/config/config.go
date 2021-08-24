package config

import (
	"fmt"
	"github.com/ktsstudio/selectel-exporter/pkg/apperrors"
	"os"
)

type ExporterConfig struct {
	Token   string
	Region  string
}

var AvailableRegions = []string{"ru-1", "ru-2", "ru-7", "ru-3", "ru-9", "ru-8"}

func Parse() (*ExporterConfig, error) {
	conf := &ExporterConfig{}

	token, ok := os.LookupEnv("SELECTEL_TOKEN")
	if !ok {
		return nil, apperrors.NewConfigError("env variable SELECTEL_TOKEN is required")
	}
	conf.Token = token

	region, ok := os.LookupEnv("SELECTEL_REGION")
	if !ok {
		return nil, apperrors.NewConfigError("env variable SELECTEL_REGION is required")
	}

	for _, v := range AvailableRegions {
		if region == v {
			conf.Region = region
		}
	}
	if conf.Region == "" {
		return nil, apperrors.NewConfigError(fmt.Sprintf("region %s is not available", region))
	}

	return conf, nil
}
