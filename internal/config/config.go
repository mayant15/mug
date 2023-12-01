package config

import (
	"errors"
	"log"
	"os"
)

type FConfig struct {
  MugHome string
  MugPackageDir string
  UserHomeDir string
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

  if err := ensureDir(mugHome); err != nil {
    return nil, err
  }
  if err := ensureDir(mugPackageDir); err != nil {
    return nil, err
  }

  _config = &FConfig{
    MugHome: mugHome,
    MugPackageDir: mugPackageDir,
    UserHomeDir: userHomeDir,
  }

  return _config, nil
}

func GetConfig() *FConfig {
  return _config
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println("Failed to create directory: ")
			return err
		}
	}
	return nil
}

