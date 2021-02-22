package cmd

// https://github.com/jessfraz/dockfmt/blob/master/format.go

import (
	"fmt"

	"github.com/spf13/cobra"
)

type cmdFmt struct {
	global *cmdGlobal
}

func (c *cmdFmt) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "fmt",
		Short: "Format the Dropletfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("fmt called")

			return nil
		},
	}
}

func init() {
	fmtCmd := cmdFmt{global: &globalCmd}
	rootCmd.AddCommand(fmtCmd.Command())
}
