package registry

import (
	"log"
)

type FPackage struct {
	Name     string    `json:"name"`
	Repo     string    `json:"repo"`
	Artifact FArtifact `json:"artifact"`
}

func (pkg FPackage) GetLatestVersionString() (string, error) {
	// TODO
	return "0.40.2", nil
	// return "14.0.3", nil
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

	err = pkg.Artifact.Fetch(fetchUrl, destDir)
	if err != nil {
		log.Println("Failed to fetch artifact: ")
		return err
	}

	return nil
}
