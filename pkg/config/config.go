package config

import (
	"errors"
	"gopkg.in/gcfg.v1"
	"strings"
)

type ServerStruct struct {
	Env  string
	Port int
}

type (
	MainConfig struct {
		Server ServerStruct
	}
)

var (
	MainCfg *MainConfig
)

func NewMainConfig(filePath string) error {

	var customErr error

	if MainCfg == nil {
		var mc MainConfig
		err := gcfg.ReadFileInto(&mc, filePath)
		if err != nil {
			errDesc := err.Error()
			customErr = errors.New("Could not load config file: " + filePath + errDesc)
			if !strings.Contains(errDesc, "can't store data at section") {
				return customErr
			}
		}

		if err != nil {
			return err
		}
		MainCfg = &mc
	}
	return customErr
}

var GetConfig = func() (*MainConfig, error) {
	if MainCfg == nil {
		return nil, errors.New("Config is empty")
	}
	return MainCfg, nil
}
