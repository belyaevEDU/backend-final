package config

import (
	"errors"
	"fmt"
)

type Config struct {
	App          AppConfig
	Philharmonic PhilharmonicConfig
}

func Load() (Config, error) {
	var errs []error

	appCfg, err := loadAppConfig()
	if err != nil {
		errs = append(errs, fmt.Errorf("loading app config: %w", err))
	}

	phrmCfg, err := loadPhilharmonicConfig()
	if err != nil {
		errs = append(errs, fmt.Errorf("loading philharmonic config: %w", err))
	}

	if len(errs) == 0 {
		return Config{App: appCfg, Philharmonic: phrmCfg}, nil
	} else {
		return Config{}, errors.Join(err)
	}
}
