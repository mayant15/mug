package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

type FPackage struct {
	Name     string    `json:"name"`
	Repo     string    `json:"repo"`
	Artifact FArtifact `json:"artifact"`
}

const (
	GITHUB_RELEASES_URL_TEMPLATE = "https://api.github.com/repos/{{ .UserName }}/{{ .RepoName }}/releases/latest"
)

func (pkg FPackage) GetLatestVersionString() (string, error) {
	userName, repoName := extractNames(pkg.Repo)
	latestReleaseUrl, err := buildLatestReleaseUrl(fReleaseUrlBuilderParams{
		UserName: userName,
		RepoName: repoName,
	})
	if err != nil {
		log.Println("Failed to build release url:")
		return "", err
	}

	response, err := http.Get(latestReleaseUrl)
	if err != nil || response.StatusCode != 200 {
		log.Println("Failed to get latest release details:")
		return "", err
	}
	defer response.Body.Close()

	var releaseInfo fReleaseInfo
	err = json.NewDecoder(response.Body).Decode(&releaseInfo)
	if err != nil {
		log.Println("Failed to parse release response body:")
		return "", err
	}

	version, err := extractVersion(releaseInfo.TagName)
	if err != nil {
		log.Println("Failed to extract version number from tag")
		return "", err
	}

	return version, nil
}

func (pkg FPackage) FetchLatestArtifact(destDir string) error {
	version, err := pkg.GetLatestVersionString()
	if err != nil {
		log.Println("Failed to fetch latest version string: ")
		return err
	}

	fetchUrl, err := pkg.Artifact.BuildUrl(FUrlBuilderParams{
		Version: version,
	})
	if err != nil {
		log.Println("Failed to build download url: ")
		return err
	}

	log.Printf("Downloading %s v%s...", pkg.Name, version)

	err = pkg.Artifact.Fetch(fetchUrl, destDir)
	if err != nil {
		log.Println("Failed to fetch artifact: ")
		return err
	}

	return nil
}

type fReleaseUrlBuilderParams struct {
	UserName string
	RepoName string
}

type fReleaseInfo struct {
	TagName string `json:"tag_name"`
}

func extractNames(repo string) (userName, repoName string) {
	if strings.HasPrefix(repo, "https://github.com") {
		segments := strings.Split(repo, "/")
		repoName = segments[len(segments)-1]
		userName = segments[len(segments)-2]
		return
	} else {
		panic("unimplemented")
	}
}

func extractVersion(tag string) (string, error) {
	regex, err := regexp.Compile(`[0-9]+\.[0-9]+\.[0-9]+`)
	if err != nil {
		log.Println("Failed to compile version regex:")
		return "", err
	}

	version := regex.FindString(tag)
	if version == "" {
		return "", errors.New("empty match")
	}

	return version, nil
}

func buildLatestReleaseUrl(params fReleaseUrlBuilderParams) (string, error) {
	templ, err := template.New("releaseUrl").Parse(GITHUB_RELEASES_URL_TEMPLATE)
	if err != nil {
		log.Println("Failed to parse release url template:")
		return "", err
	}

	var buf bytes.Buffer
	err = templ.Execute(&buf, params)
	if err != nil {
		log.Println("Failed to build release url template:")
		return "", err
	}

	return buf.String(), nil
}
