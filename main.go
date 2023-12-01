package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

type EArtifactType int

const (
  eArtifactType_Tarball EArtifactType = 0
)

type FUrlBuilderParams struct {
  Version string
}

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

func (artifact FArtifact) BuildUrl(params FUrlBuilderParams) (string, error) {
  templ, err := template.New("artifactUrl").Parse(artifact.Url)
  if err != nil {
    log.Println("Failed to parse url template: ")
    return "", err
  }

  var buf bytes.Buffer
  err = templ.Execute(&buf, params)
  if err != nil {
    log.Println("Failed to build url template: ")
    return "", err
  }

  return buf.String(), nil
}

func downloadTarball(fileUrl string, destDir string) error {
  log.Println("downloading tarball...")

  parsedUrl, err := url.Parse(fileUrl)
  if err != nil {
    log.Println("Failed to parse fileUrl: ")
    return err
  }

  segments := strings.Split(parsedUrl.Path, "/")
  filename := segments[len(segments) - 1]

  destPath := destDir + "/" + filename

  destFile, err := os.Create(destPath)
  if err != nil {
    log.Println("Failed to create file: ")
    return err
  }

  client := http.Client{
    CheckRedirect: func(r *http.Request, via []*http.Request) error {
      r.URL.Opaque = r.URL.Path
      return nil
    },
  }

  response, err := client.Get(fileUrl)
  if err != nil {
    log.Println("Failed to download file: ")
    return err
  }
  defer response.Body.Close()

  size, err := io.Copy(destFile, response.Body)
  defer destFile.Close()
  
  fmt.Printf("Downloaded file %s of size %d", filename, size)
  return nil
}

func (artifact FArtifact) Fetch(url string, destDir string) error {
  switch artifact.ArtifactType {
  case eArtifactType_Tarball:
    return downloadTarball(url, destDir)
  }
  return nil
}

func (reg FRegistry) FindPackage(name string) (*FPackage, error) {
  for i := range reg.Pkgs {
    if reg.Pkgs[i].Name == name {
      return &reg.Pkgs[i], nil
    }
  }
  return nil, errors.New("package not found in registry")
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

func innerMain() error {
  userHomeDir, err := os.UserHomeDir()
  if err != nil {
    log.Println("Failed to get user's home directory")
    return err
  }

  MUG_HOME := userHomeDir + "/.mug"
  MUG_PACKAGES_DIR := MUG_HOME + "/packages"

  if err := ensureDir(MUG_HOME); err != nil { return err }
  if err := ensureDir(MUG_PACKAGES_DIR); err != nil { return err }

  registry, err := loadRegistry()
  if err != nil {
    return err
  }

  rg, err := registry.FindPackage("ripgrep")
  if err != nil {
    log.Println("Failed to locate package: ")
    return err
  }

  artifact := rg.Artifact
  fetchUrl, err := artifact.BuildUrl(FUrlBuilderParams{
    Version: "14.0.3",
  })
  artifact.Fetch(fetchUrl, MUG_PACKAGES_DIR)

  return nil
}

func main() {
  err := innerMain()
  if err != nil {
    log.Fatal("\t", err)
  }
}
