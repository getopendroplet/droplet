package builder

import (
	"fmt"

	"github.com/getopendroplet/droplet/config"
	"github.com/getopendroplet/droplet/dropletfile/instructions"
)

// Builder - interface
type Builder interface {
	Arg(command instructions.ArgCommand) string
	Config(command instructions.ConfigCommand) string
	Copy(command instructions.CopyCommand) string
	Cron(command instructions.CronCommand) string
	Delete(command instructions.DeleteCommand) string
	Env(command instructions.EnvCommand) string
	Expose(command instructions.ExposeCommand) string
	Label(command instructions.LabelCommand) string
	Run(command instructions.RunCommand) string
	User(command instructions.UserCommand) string
	Package(command instructions.PackageCommand) string
	Workdir(command instructions.WorkdirCommand) string
}

// Build - Dropletfile to script
func Build(conf *config.Config, commands []instructions.Command) error {
	b := LocalBuilder{conf: conf}

	for _, c := range commands {
		fmt.Printf("Building command: %s\n", c.Name())
		result := ""
		switch cmd := c.(type) {
		case *instructions.ArgCommand:
			result = b.Arg(*cmd)
		case *instructions.ConfigCommand:
			result = b.Config(*cmd)
		case *instructions.CopyCommand:
			result = b.Copy(*cmd)
		case *instructions.CronCommand:
			result = b.Cron(*cmd)
		case *instructions.DeleteCommand:
			result = b.Delete(*cmd)
		case *instructions.EnvCommand:
			result = b.Env(*cmd)
		case *instructions.ExposeCommand:
			result = b.Expose(*cmd)
		case *instructions.LabelCommand:
			result = b.Label(*cmd)
		case *instructions.RunCommand:
			result = b.Run(*cmd)
		case *instructions.UserCommand:
			result = b.User(*cmd)
		case *instructions.PackageCommand:
			result = b.Package(*cmd)
		case *instructions.WorkdirCommand:
			result = b.Workdir(*cmd)
		}

		fmt.Printf("Result: %s\n", result)
	}
	return nil
}
