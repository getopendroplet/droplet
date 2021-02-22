package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/getopendroplet/droplet/utils"
	"github.com/getopendroplet/droplet/utils/table"

	"github.com/spf13/cobra"
)

type cmdConfig struct {
	global *cmdGlobal
}

type cmdConfigSet struct {
	global *cmdGlobal
}

type cmdConfigGet struct {
	global *cmdGlobal
}

type cmdConfigDel struct {
	global *cmdGlobal
}

type cmdConfigList struct {
	global *cmdGlobal

	flagFormat string
}

func (c *cmdConfig) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage Droplet configurations",
	}
}

func (c *cmdConfigSet) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key>=<value> ...",
		Short: "Set a configuration key",
		Args:  cobra.MinimumNArgs(1),
		RunE:  c.Run,
	}
}

func (c *cmdConfigSet) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf

	for _, arg := range args {
		fields := strings.SplitN(arg, "=", 2)
		if len(fields) != 2 {
			return fmt.Errorf("Invalid key=value configuration: %s", arg)
		}

		if _, ok := conf.Configs[fields[0]]; ok {
			return fmt.Errorf("Config %s already exists", fields[0])
		}

		conf.Configs[fields[0]] = fields[1]
	}

	return nil
}

func (c *cmdConfigGet) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE:  c.Run,
	}
}

func (c *cmdConfigGet) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf
	key := args[0]

	if _, ok := conf.Configs[key]; !ok {
		return fmt.Errorf("Config %s doesn't exist", key)
	}

	fmt.Println(conf.Configs[key])

	return nil
}

func (c *cmdConfigDel) Command() *cobra.Command {
	return &cobra.Command{
		Use:     "del <key>",
		Aliases: []string{"rm"},
		Short:   "Delete a configuration key",
		Args:    cobra.ExactArgs(1),
		RunE:    c.Run,
	}
}

func (c *cmdConfigDel) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf
	key := args[0]

	if _, ok := conf.Configs[key]; !ok {
		return fmt.Errorf("Config %s doesn't exist", key)
	}

	delete(conf.Configs, args[0])

	return nil
}

func (c *cmdConfigList) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List configurations",
		Args:    cobra.ExactArgs(0),
		RunE:    c.Run,
	}

	cmd.Flags().StringVar(&c.flagFormat, "format", "table", "Format (csv|json|table|yaml)")
	return cmd
}

func (c *cmdConfigList) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf
	// List the aliases
	data := [][]string{}
	for k, v := range conf.Configs {
		data = append(data, []string{k, v})
	}
	sort.Sort(utils.ByName(data))

	header := []string{"Key", "Value"}
	return table.RenderTable(c.flagFormat, header, data, conf.Configs)
}

func init() {
	configCmd := cmdConfig{global: &globalCmd}
	configSetCmd := cmdConfigSet{global: &globalCmd}
	configGetCmd := cmdConfigGet{global: &globalCmd}
	configDelCmd := cmdConfigDel{global: &globalCmd}
	configListCmd := cmdConfigList{global: &globalCmd}

	config := configCmd.Command()
	config.AddCommand(configSetCmd.Command())
	config.AddCommand(configGetCmd.Command())
	config.AddCommand(configDelCmd.Command())
	config.AddCommand(configListCmd.Command())

	rootCmd.AddCommand(config)
}
