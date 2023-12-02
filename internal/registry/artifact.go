package registry

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mayant15/mug/internal/util"
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
	Url          string        `json:"url"`
	BinaryPath   string        `json:"binaryPath"`
	filePath     string
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

func (artifact *FArtifact) Fetch(url string, destDir string) error {
	switch artifact.ArtifactType {
	case eArtifactType_Tarball:
		path, err := downloadTarball(url, destDir)
		artifact.filePath = path
		return err

	default:
		log.Println("Failed to fetch artifact: ")
		return errors.New("unsupported artifact type")
	}
}

func (artifact *FArtifact) Prepare() error {
	return nil
}

func (artifact FArtifact) Install() error {
	return nil
}

func downloadTarball(fileUrl string, destDir string) (string, error) {
	err := util.EnsureDir(destDir)
	if err != nil {
		log.Println("Failed to create download directory:")
		return "", err
	}

	parsedUrl, err := url.Parse(fileUrl)
	if err != nil {
		log.Println("Failed to parse fileUrl: ")
		return "", err
	}

	segments := strings.Split(parsedUrl.Path, "/")
	filename := segments[len(segments)-1]

	destPath := destDir + "/" + filename

	destFile, err := os.Create(destPath)
	if err != nil {
		log.Println("Failed to create file: ")
		return "", err
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
		return "", err
	}

	if response.StatusCode == 404 {
		log.Println("Failed to download file: ")
		return "", errors.New("404 not found")
	}

	defer response.Body.Close()

	size, err := io.Copy(destFile, response.Body)
	defer destFile.Close()

	log.Printf("Downloaded file %s of size %d", filename, size)
	return destPath, nil
}
