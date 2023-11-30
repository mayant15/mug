package main

import (
	"os"
	"fmt"
  "encoding/json"
  "log"
)

type EArtifactType int

const (
  TARBALL EArtifactType = 0
)

type FArtifact struct {
  ArtifactType EArtifactType `json:"type"`
  Url string `json:"url"`
  BinaryPath string `json:"binaryPath"`
}

type FPackage struct {
  Name string `json:"name"`
  Repo string `json:"repo"`
  Artifact FArtifact `json:"artifact"`
}

type FRegistry struct {
  Pkgs []FPackage `json:"pkgs"`
}

func loadRegistry() (*FRegistry, error) {
  bytes, err := os.ReadFile("./resources/registry.json")
  if err != nil {
    log.Println("Failed to read registry file:")
    return nil, err
  }

  var registry FRegistry
  err = json.Unmarshal(bytes, &registry)
  if err != nil {
    log.Println("Failed to parse registry JSON:")
    return nil, err
  }

  return &registry, nil
}

func main() {
  registry, err := loadRegistry()
  if err != nil {
    log.Fatal("\t", err)
  }

  fmt.Println(registry.Pkgs[0].Artifact.Url)
}
