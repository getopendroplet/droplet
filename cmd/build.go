package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/getopendroplet/droplet/dropletfile/instructions"
	"github.com/getopendroplet/droplet/dropletfile/parser"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type cmdBuild struct {
	global *cmdGlobal
}

func (c *cmdBuild) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "build <dropletfile> <stage>",
		Short: "Build an script from a Dropletfile",
		Args:  cobra.ExactArgs(2),
		RunE:  c.Run,
	}
}

func (c *cmdBuild) Run(cmd *cobra.Command, args []string) error {
	// conf := c.global.conf
	fileName := args[0]
	stageName := args[1]

	var f *os.File
	var err error

	f, err = os.Open(fileName + "/Dropletfile")
	if err != nil {
		return err
	}
	defer f.Close()

	result, err := parser.Parse(f)
	if err != nil {
		return err
	}

	stages, _, err := instructions.Parse(result.AST)
	if err != nil {
		return err
	}

	index, exists := instructions.HasStage(stages, stageName)
	if !exists {
		return errors.Errorf("no build stage %s in current Dropletfile", stageName)
	}

	stage := stages[index]

	json, err := json.Marshal(stage)
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}

func init() {
	buildCmd := cmdBuild{global: &globalCmd}
	rootCmd.AddCommand(buildCmd.Command())
}
