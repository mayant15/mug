package cmd

import (
	"log"

	"github.com/mayant15/mug/internal/config"
	"github.com/mayant15/mug/internal/registry"
	"github.com/mayant15/mug/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install <package>",
	Short: "Install a package from the registry",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.InitConfig()
		if err != nil {
			return err
		}

		return handleInstallCmd(args)
	},
}

func handleInstallCmd(packages []string) error {
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

		if checkInstalled(*pkg, config.MugInstallDir) {
			log.Printf("Package %s is already installed", pkgName)
			continue
		}

		err = doInstall(pkg, config.MugPackageDir)
		if err != nil {
			log.Printf("!!! FAILED TO INSTALL %s !!!", pkgName)
			log.Println("    ERROR: ", err)
		}
	}

	return err
}

func doInstall(pkg *registry.FPackage, packagesDir string) error {
	err := pkg.FetchLatestArtifact(packagesDir)
	if err != nil {
		log.Println("Failed to fetch latest artifact for package: ")
		return err
	}

	err = pkg.Prepare()
	if err != nil {
		log.Println("Failed to prepare artifact for installation: ")
		return err
	}

	err = pkg.Install()
	if err != nil {
		log.Println("Failed to install artifact:")
		return err
	}

	return nil
}

func checkInstalled(pkg registry.FPackage, installDir string) bool {
	link := pkg.Artifact.BuildSymLinkPath(installDir)
	return util.CheckExists(link)
}
