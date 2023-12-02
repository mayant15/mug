package registry

import (
	"bytes"
	"errors"
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
	Url          string        `json:"url"`
	BinaryPath   string        `json:"binaryPath"`
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
	parsedUrl, err := url.Parse(fileUrl)
	if err != nil {
		log.Println("Failed to parse fileUrl: ")
		return err
	}

	segments := strings.Split(parsedUrl.Path, "/")
	filename := segments[len(segments)-1]

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

	if response.StatusCode == 404 {
		log.Println("Failed to download file: ")
		return errors.New("404 not found")
	}

	defer response.Body.Close()

	size, err := io.Copy(destFile, response.Body)
	defer destFile.Close()

	log.Printf("Downloaded file %s of size %d", filename, size)
	return nil
}

func (artifact FArtifact) Fetch(url string, destDir string) error {
	switch artifact.ArtifactType {
	case eArtifactType_Tarball:
		return downloadTarball(url, destDir)

	default:
		log.Println("Failed to fetch artifact: ")
		return errors.New("unsupported artifact type")
	}
}
