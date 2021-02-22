// Package command contains the set of Dropletfile commands.
package command

// Define constants for the command strings
const (
	Arg     = "arg"
	Config  = "config"
	Copy    = "copy"
	Cron    = "cron"
	Delete  = "delete"
	Env     = "env"
	Expose  = "expose"
	Label   = "label"
	Package = "package"
	Run     = "run"
	Stage   = "stage"
	User    = "user"
	Workdir = "workdir"
)

// Commands is list of all Dropletfile commands
var Commands = map[string]struct{}{
	Arg:     {},
	Config:  {},
	Copy:    {},
	Cron:    {},
	Delete:  {},
	Env:     {},
	Expose:  {},
	Label:   {},
	Package: {},
	Run:     {},
	Stage:   {},
	User:    {},
	Workdir: {},
}
