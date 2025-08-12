package config

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
)

func GetConfigDir() string {
	return filepath.Join(xdg.ConfigHome, APPNAME)
}

func GetConfigFilePath() (string, error) {
	cfgPath, err := xdg.ConfigFile(filepath.Join(APPNAME, "config.yml"))
	if err != nil {
		return "", fmt.Errorf("error while getting the configuration file path")
	}
	return cfgPath, nil
}
