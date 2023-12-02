package registry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mayant15/mug/internal/config"
	"github.com/mayant15/mug/internal/util"
	"github.com/ulikunitz/xz"
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
	BinaryAlias  string        `json:"alias"`
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

func (artifact *FArtifact) Install() error {
	switch artifact.ArtifactType {
	case eArtifactType_Tarball:
		return installTarball(*artifact)

	default:
		log.Println("Failed to install artifact:")
		return errors.New("unsupported artifact type")
	}
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

	_, err = io.Copy(destFile, response.Body)
	defer destFile.Close()

	return destPath, nil
}

func installTarball(artifact FArtifact) error {
	destDir, tarfile := path.Split(artifact.filePath)
	err := extractTarball(artifact.filePath, destDir)
	if err != nil {
		log.Println("Failed to extract tarball:")
		return err
	}

	extIdx := strings.Index(tarfile, ".tar")
	if extIdx == -1 {
		log.Fatalln("Tar archives must have a .tar extension")
	}

	tarname := tarfile[:extIdx]

	binary, err := findExtractedBinary(destDir, tarname, artifact.BinaryPath)
	if err != nil {
		log.Printf("Failed to find binary %s:", artifact.BinaryPath)
		return err
	}

	err = linkBinary(binary, artifact.BinaryAlias)
	if err != nil {
		log.Println("Failed to symlink binary:")
		return err
	}

	return nil
}

func linkBinary(file string, alias string) error {
	installDir := config.GetConfig().MugInstallDir

	var name = alias
	if name == "" {
		name = path.Base(file)
	}

	newname := path.Join(installDir, name)
	return os.Symlink(file, newname)
}

func findExtractedBinary(dir string, tarname string, binaryPath string) (string, error) {
	attempt := path.Join(dir, binaryPath)
	if util.CheckExists(attempt) {
		return path.Clean(attempt), nil
	}

	attempt = path.Join(dir, tarname, binaryPath)
	if util.CheckExists(attempt) {
		return path.Clean(attempt), nil
	}

	return "", errors.New("could not find binary")
}

func extractTarball(tarpath string, destDir string) error {
	tarfile, err := openTarfile(tarpath)
	if err != nil {
		log.Println("Failed to open downloaded artifact:")
		return err
	}

	reader := tar.NewReader(tarfile)

	header, err := reader.Next()
	for err != io.EOF {

		if err != nil {
			log.Println("Failed to read tarfile:")
			return err
		}

		filepath := path.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filepath, os.ModePerm)
			if err != nil {
				log.Printf("Failed to create directory %s:", header.Name)
				return err
			}

		case tar.TypeReg:
			file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				log.Printf("Failed to create file %s:", header.Name)
				return err
			}

			_, err = io.Copy(file, reader)
			if err != nil {
				log.Println("Failed to copy file contents from tarball:")
				return err
			}

			file.Close()

		default:
			log.Println("Failed to extract tarball:")
			return errors.New("unknown type flag")
		}

		header, err = reader.Next()
	}

	return nil
}

func openTarfile(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Failed to read downloaded artifact:")
		return nil, err
	}

	if strings.HasSuffix(path, ".tar.gz") {
		tarfile, err := gzip.NewReader(file)
		if err != nil {
			log.Println("Failed to uncompress .tar.gz:")
			return nil, err
		}
		return tarfile, nil
	} else if strings.HasSuffix(path, ".tar.xz") {
		tarfile, err := xz.NewReader(file)
		if err != nil {
			log.Println("Failed to uncompress .tar.xz:")
			return nil, err
		}
		return tarfile, nil
	} else if strings.HasSuffix(path, ".tar") {
		return file, nil
	}

	return nil, errors.New("unsupported filetype: must be one of .tar, .tar.gz")
}
