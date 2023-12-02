package util

import (
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
