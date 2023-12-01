package registry

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type FRegistry struct {
	Pkgs []FPackage `json:"pkgs"`
}

func LoadRegistryFromFile() (*FRegistry, error) {
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

func (reg FRegistry) FindPackage(name string) (*FPackage, error) {
	for i := range reg.Pkgs {
		if reg.Pkgs[i].Name == name {
			return &reg.Pkgs[i], nil
		}
	}
	return nil, errors.New("package not found in registry")
}
