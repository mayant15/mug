package cmd

import (
	"log"

	"github.com/mayant15/mug/internal/config"
	"github.com/mayant15/mug/internal/registry"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update <package>",
	Short: "Update an already installed package",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.InitConfig()
		if err != nil {
			return err
		}

		return handleUpdateCmd(args)
	},
}

func handleUpdateCmd(packages []string) error {
	conf := config.GetConfig()

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

		if !checkInstalled(*pkg, conf.MugInstallDir) {
			log.Printf("Package %s is not installed", pkgName)
			continue
		}

		err = doUpdate(pkg, *conf)
		if err != nil {
			log.Printf("!!! FAILED TO UPDATE %s !!!", pkgName)
			log.Println("    ERROR: ", err)
		}
	}

	return err
}

func doUpdate(pkg *registry.FPackage, conf config.FConfig) error {
	err := doRemove(pkg, conf.MugInstallDir)
	if err != nil {
		return err
	}
	return doInstall(pkg, conf.MugPackageDir)
}
