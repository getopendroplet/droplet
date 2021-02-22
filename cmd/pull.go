package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type cmdPull struct {
	global *cmdGlobal
}

func (c *cmdPull) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "pull <remote> <droplet>",
		Short: "Pull an Droplet from a remote",
		Args: func(cmd *cobra.Command, args []string) error {
			n := 2
			if len(args) != n {
				return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			remote := args[0]
			droplet := args[1]
			fmt.Println("pull from " + remote + " an " + droplet)

			return nil
		},
	}
}

func init() {
	pullCmd := cmdPull{global: &globalCmd}
	rootCmd.AddCommand(pullCmd.Command())
}
