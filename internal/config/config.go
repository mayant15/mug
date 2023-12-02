package config

import (
	"errors"
	"log"
	"os"

	"github.com/mayant15/mug/internal/util"
)

type FConfig struct {
	MugHome       string
	MugPackageDir string
	UserHomeDir   string
}

var _config *FConfig = nil

func InitConfig() (*FConfig, error) {
	if _config != nil {
		return nil, errors.New("config already initialized")
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Failed to get user's home directory")
		return nil, err
	}

	mugHome := userHomeDir + "/.mug"
	mugPackageDir := mugHome + "/packages"

	if err := util.EnsureDir(mugHome); err != nil {
		return nil, err
	}
	if err := util.EnsureDir(mugPackageDir); err != nil {
		return nil, err
	}

	_config = &FConfig{
		MugHome:       mugHome,
		MugPackageDir: mugPackageDir,
		UserHomeDir:   userHomeDir,
	}

	return _config, nil
}

func GetConfig() *FConfig {
	return _config
}
