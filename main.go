package main

import (
	"log"
	"os"

  "mayant15/mug/internal/registry"
)

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

func innerMain() error {
  const packageName = "lazygit"

  userHomeDir, err := os.UserHomeDir()
  if err != nil {
    log.Println("Failed to get user's home directory")
    return err
  }

  MUG_HOME := userHomeDir + "/.mug"
  MUG_PACKAGES_DIR := MUG_HOME + "/packages"

  if err := ensureDir(MUG_HOME); err != nil { return err }
  if err := ensureDir(MUG_PACKAGES_DIR); err != nil { return err }

  registry, err := registry.LoadRegistryFromFile()
  if err != nil {
    return err
  }

  pkg, err := registry.FindPackage(packageName)
  if err != nil {
    log.Println("Failed to locate package: ")
    return err
  }

  err = pkg.FetchLatestArtifact(MUG_PACKAGES_DIR)
  if err != nil {
    log.Println("Failed to fetch latest artifact for package: ")
    return err
  }
  
  return nil
}

func main() {
  err := innerMain()
  if err != nil {
    log.Fatal("\t", err)
  }
}
