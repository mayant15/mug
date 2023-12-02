package util

import (
	"errors"
	"log"
	"os"
)

func EnsureDir(path string) error {
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

func CheckExists(path string) bool {
  _, err := os.Stat(path)
  return err == nil
}

func CreateFile(path string) (*os.File, error) {
  if CheckExists(path) {
    log.Printf("Failed to create file %s: ", path)
    return nil, errors.New("file already exists")
  }

  file, err := os.Create(path)
  if err != nil {
    log.Printf("Failed to create file %s: ", path)
    return nil, err
  }

  return file, nil
}

func EnsureFile(path string) error {
  if CheckExists(path) {
    return nil
  }

  _, err := os.Create(path)
  if err != nil {
    log.Printf("Failed to create file %s: ", path)
    return err
  }

  return nil
}
