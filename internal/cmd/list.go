package cmd

import (
	"strings"

	"github.com/mayant15/mug/internal/config"
	"github.com/mayant15/mug/internal/registry"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.InitConfig()
		if err != nil {
			return err
		}

		return handleListCmd()
	},
}

func handleListCmd() error {
	conf := config.GetConfig()

	reg, err := registry.LoadRegistryFromFile()
	if err != nil {
		return err
	}

	names := []string{}
	for _, pkg := range reg.Pkgs {
		if checkInstalled(pkg, conf.MugInstallDir) {
			names = append(names, "* "+pkg.Name)
		}
	}

	println(strings.Join(names, "\n"))
	return err
}
