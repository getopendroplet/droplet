package cmd

import (
	"fmt"
	"sort"

	"github.com/getopendroplet/droplet/config"
	"github.com/getopendroplet/droplet/utils"
	"github.com/getopendroplet/droplet/utils/table"

	"github.com/spf13/cobra"
)

type cmdRemote struct {
	global *cmdGlobal
}

type cmdRemoteSet struct {
	global *cmdGlobal
}

type cmdRemoteDel struct {
	global *cmdGlobal
}

type cmdRemoteList struct {
	global *cmdGlobal

	flagFormat string
}

type cmdRemoteSwitch struct {
	global *cmdGlobal
}

func (c *cmdRemote) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "remote",
		Short: "Manage Droplet the list of remotes",
	}
}

func (c *cmdRemoteSet) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <address> <protocol>",
		Short: "Set a new remote",
		Args:  cobra.ExactArgs(3),
		RunE:  c.Run,
	}
}

func (c *cmdRemoteSet) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf
	name := args[0]
	addr := args[1]
	protocol := args[2]

	if _, ok := conf.Remotes[name]; ok {
		return fmt.Errorf("Remote %s already exists", name)
	}

	conf.Remotes[name] = config.Remote{
		Addr:     addr,
		Protocol: protocol,
	}

	return nil
}

func (c *cmdRemoteDel) Command() *cobra.Command {
	return &cobra.Command{
		Use:     "del <name>",
		Aliases: []string{"rm"},
		Short:   "Del a remote",
		Args:    cobra.ExactArgs(1),
		RunE:    c.Run,
	}
}

func (c *cmdRemoteDel) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf
	name := args[0]

	if _, ok := conf.Remotes[name]; !ok {
		return fmt.Errorf("Remote %s doesn't exist", name)
	}

	if conf.DefaultRemote == name {
		return fmt.Errorf("Can't delete the default remote")
	}

	delete(conf.Remotes, name)

	return nil
}

func (c *cmdRemoteList) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the available remotes",
		Args:    cobra.ExactArgs(0),
		RunE:    c.Run,
	}

	cmd.Flags().StringVar(&c.flagFormat, "format", "table", "Format (csv|json|table|yaml)")
	return cmd
}

func (c *cmdRemoteList) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf

	data := [][]string{}
	for name, rc := range conf.Remotes {
		strName := name
		if name == conf.DefaultRemote {
			strName = fmt.Sprintf("%s (%s)", name, "default")
		}
		data = append(data, []string{strName, rc.Addr, rc.Protocol})
	}
	sort.Sort(utils.ByName(data))

	header := []string{"Name", "Url", "Protocol"}
	return table.RenderTable(c.flagFormat, header, data, conf.Configs)
}

func (c *cmdRemoteSwitch) Command() *cobra.Command {
	return &cobra.Command{
		Use:     "switch <name>",
		Aliases: []string{"set-default"},
		Short:   "Switch the default remote",
		Args:    cobra.ExactArgs(1),
		RunE:    c.Run,
	}
}

func (c *cmdRemoteSwitch) Run(cmd *cobra.Command, args []string) error {
	conf := c.global.conf
	name := args[0]

	if _, ok := conf.Remotes[name]; !ok {
		return fmt.Errorf("Config %s doesn't exist", name)
	}

	conf.DefaultRemote = name

	return nil
}

func init() {
	remoteCmd := cmdRemote{global: &globalCmd}
	remoteSetCmd := cmdRemoteSet{global: &globalCmd}
	remoteDelCmd := cmdRemoteDel{global: &globalCmd}
	remoteListCmd := cmdRemoteList{global: &globalCmd}
	remoteSwitchCmd := cmdRemoteSwitch{global: &globalCmd}

	remote := remoteCmd.Command()
	remote.AddCommand(remoteSetCmd.Command())
	remote.AddCommand(remoteDelCmd.Command())
	remote.AddCommand(remoteListCmd.Command())
	remote.AddCommand(remoteSwitchCmd.Command())

	rootCmd.AddCommand(remote)
}
