package cmd

import (
	"droplet/config"
	"droplet/utils"
	"droplet/utils/stack"
	"droplet/version"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

type cmdGlobal struct {
	conf          *config.Config
	confPath      string
	workspacePath string
	cmd           *cobra.Command
	ret           int

	flagHelp       bool
	flagLogDebug   bool
	flagLogVerbose bool
	flagQuiet      bool
}

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "droplet",
		Short: "Droplet is a bash script builder",
		Long: `Droplet is a bash script builder
	
All of Droplet's features can be driven through the various commands below.
For help with any of those, simply call them with --help.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	globalCmd = cmdGlobal{cmd: rootCmd}
)

func init() {
	stack.SetVersionInfo(version.Version, version.Revision)

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&globalCmd.workspacePath, "workspace", home+"/droplet", "Set workspace path")
	rootCmd.PersistentFlags().StringVar(&globalCmd.confPath, "config", home+"/.config/droplet", "Set config path")
	rootCmd.PersistentFlags().BoolVarP(&globalCmd.flagHelp, "help", "h", false, "Print help")
	rootCmd.PersistentFlags().BoolVar(&globalCmd.flagLogDebug, "debug", false, "Show all debug messages")
	rootCmd.PersistentFlags().BoolVarP(&globalCmd.flagLogVerbose, "verbose", "v", false, "Show all information messages")
	rootCmd.PersistentFlags().BoolVarP(&globalCmd.flagQuiet, "quiet", "q", false, "Don't show progress information")

	// Wrappers
	rootCmd.PersistentPreRunE = globalCmd.PreRun
	rootCmd.PersistentPostRunE = globalCmd.PostRun
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if globalCmd.flagLogDebug {
			fmt.Fprintf(os.Stderr, "error: %+v", stack.Formatter(err))
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(1)
	}
}

func (c *cmdGlobal) PreRun(cmd *cobra.Command, args []string) error {
	var err error
	configFile := os.ExpandEnv(path.Join(c.confPath, "config.yml"))

	if c.conf, err = config.LoadConfig(configFile); err != nil {
		c.conf = config.NewConfig(true)
	}

	// Create the workspace dir so that we don't get in here again for this user.
	if !utils.PathExists(c.workspacePath) {
		if err = os.MkdirAll(c.workspacePath, 0750); err != nil {
			return err
		}
	}

	return nil
}

func (c *cmdGlobal) PostRun(cmd *cobra.Command, args []string) error {
	configFile := os.ExpandEnv(path.Join(c.confPath, "config.yml"))

	if err := c.conf.SaveConfig(configFile); err != nil {
		return err
	}

	return nil
}
