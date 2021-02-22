package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type cmdVersion struct {
	global *cmdGlobal
}

func (c *cmdVersion) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show Droplet version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("version called")

			return nil
		},
	}
}

func init() {
	versionCmd := cmdVersion{global: &globalCmd}
	rootCmd.AddCommand(versionCmd.Command())
}
