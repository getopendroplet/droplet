package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type cmdPush struct {
	global *cmdGlobal
}

func (c *cmdPush) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "push <remote> <droplet>",
		Short: "Push an Droplet to a remote",
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
			fmt.Println("push to " + remote + " an " + droplet)

			return nil
		},
	}
}

func init() {
	pushCmd := cmdPush{global: &globalCmd}
	rootCmd.AddCommand(pushCmd.Command())
}
