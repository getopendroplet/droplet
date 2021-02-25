package builder

import (
	"strings"

	"github.com/getopendroplet/droplet/config"
	"github.com/getopendroplet/droplet/dropletfile/instructions"
)

// LocalBuilder - build local commands
type LocalBuilder struct {
	conf *config.Config
}

// Arg - build local arg command
func (l LocalBuilder) Arg(command instructions.ArgCommand) string {
	return ""
}

// Config - build local config command
func (l LocalBuilder) Config(command instructions.ConfigCommand) string {
	return ""
}

// Copy - build local copy command
func (l LocalBuilder) Copy(command instructions.CopyCommand) string {
	cmd := []string{l.conf.Configs["builder.local.copy"]}
	return strings.Join(cmd, " ")
}

// Cron - build local cron command
func (l LocalBuilder) Cron(command instructions.CronCommand) string {
	return ""
}

// Delete - build local delete command
func (l LocalBuilder) Delete(command instructions.DeleteCommand) string {
	cmd := []string{l.conf.Configs["builder.local.delete"]}
	return strings.Join(cmd, " ")
}

// Env - build local env command
func (l LocalBuilder) Env(command instructions.EnvCommand) string {
	cmd := []string{l.conf.Configs["builder.local.env"]}
	return strings.Join(cmd, " ")
}

// Expose - build local expose command
func (l LocalBuilder) Expose(command instructions.ExposeCommand) string {
	cmd := []string{l.conf.Configs["builder.local.expose"]}
	return strings.Join(cmd, " ")
}

// Label - build local label command
func (l LocalBuilder) Label(command instructions.LabelCommand) string {
	cmd := []string{l.conf.Configs["builder.local.label"]}
	return strings.Join(cmd, " ")
}

// Run -build local run command
func (l LocalBuilder) Run(command instructions.RunCommand) string {
	return ""
}

// User -build local user command
func (l LocalBuilder) User(command instructions.UserCommand) string {
	cmd := []string{l.conf.Configs["builder.local.user"]}
	return strings.Join(cmd, " ")
}

// Package - build local package command
func (l LocalBuilder) Package(command instructions.PackageCommand) string {
	return ""
}

// Workdir - build local workdir command
func (l LocalBuilder) Workdir(command instructions.WorkdirCommand) string {
	cmd := []string{l.conf.Configs["builder.local.workdir"]}
	return strings.Join(cmd, " ")
}
