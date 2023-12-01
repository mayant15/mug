package cmd

import (
	"log"

	"github.com/mayant15/mug/internal/config"
	"github.com/mayant15/mug/internal/registry"
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

		return handleInstallCmd(args[0])
	},
}

func handleInstallCmd(pkgName string) error {
	config := config.GetConfig()
	registry, err := registry.LoadRegistryFromFile()
	if err != nil {
		return err
	}

	pkg, err := registry.FindPackage(pkgName)
	if err != nil {
		log.Println("Failed to locate package: ")
		return err
	}

	err = pkg.FetchLatestArtifact(config.MugPackageDir)
	if err != nil {
		log.Println("Failed to fetch latest artifact for package: ")
		return err
	}

	return nil
}
