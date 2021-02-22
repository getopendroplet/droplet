package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

type cmdCompletion struct {
	global *cmdGlobal
}

func (c *cmdCompletion) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generate completion script",
		Long: `To load completions:
	
	Bash:
	
	$ source <(droplet completion)
	
	# To load completions for each session, execute once:
	Linux:
	  $ droplet completion > /etc/bash_completion.d/droplet
	MacOS:
	  $ droplet completion > /usr/local/etc/bash_completion.d/droplet
	`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().GenBashCompletion(os.Stdout)

			return nil
		},
	}
}

func init() {
	completionCmd := cmdCompletion{global: &globalCmd}
	rootCmd.AddCommand(completionCmd.Command())

}
