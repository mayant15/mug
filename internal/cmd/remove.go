package cmd

import (
	"log"
	"os"

	"github.com/mayant15/mug/internal/config"
	"github.com/mayant15/mug/internal/registry"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove <package>",
	Short: "Remove an installed package",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.InitConfig()
		if err != nil {
			return err
		}

		return handleRemoveCmd(args)
	},
}

func handleRemoveCmd(packages []string) error {
	config := config.GetConfig()

	reg, err := registry.LoadRegistryFromFile()
	if err != nil {
		return err
	}

	for _, pkgName := range packages {
		pkg, err := reg.FindPackage(pkgName)
		if err != nil {
			log.Println("Failed to locate package: ")
			return err
		}

		if !checkInstalled(*pkg, config.MugInstallDir) {
			log.Printf("Package %s is not installed", pkgName)
			continue
		}

		err = doRemove(pkg, config.MugInstallDir)
		if err != nil {
			log.Printf("!!! FAILED TO REMOVE %s !!!", pkgName)
			log.Println("    ERROR: ", err)
		}
	}

	return err
}

func doRemove(pkg *registry.FPackage, installDir string) error {
	link := pkg.Artifact.BuildSymLinkPath(installDir)
	return os.Remove(link)
}
